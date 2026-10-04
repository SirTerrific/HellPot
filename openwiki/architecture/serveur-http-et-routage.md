---
type: architecture
title: Serveur HTTP, routage et piège
description: Comment HellPot configure son serveur fasthttp, décide quelles requêtes tombent dans le piège (catchall, paths, robots.txt, liste noire de user agents), identifie le client et écoute en TCP ou sur socket Unix.
tags: [http, fasthttp, routage, robots-txt, socket-unix]
sources:
  - id: openwiki-source-c34322bad84de7b6c0a25de1
    resource: repo://cmd/HellPot/HellPot.go
  - id: openwiki-source-a9ca8194716b05765ae68c77
    resource: repo://internal/http/robots.go
  - id: openwiki-source-5dcbf61c5847f06737a58a9d
    resource: repo://internal/http/router_unix_test.go
  - id: openwiki-source-fd1f4e266c537f5376b692a9
    resource: repo://internal/http/router_unix.go
  - id: openwiki-source-313641b4fcdb85c8da9945ee
    resource: repo://internal/http/router_windows.go
  - id: openwiki-source-e689a4a46f2ebef989178800
    resource: repo://internal/http/router.go
generated: { by: "claude-code", at: "2026-10-04T19:20:53.104Z" }
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T19:20:53.104Z
---

# Serveur HTTP, routage et piège

Tout le code HTTP est dans le package `internal/http` ([router.go](../../internal/http/router.go), [robots.go](../../internal/http/robots.go), [router_unix.go](../../internal/http/router_unix.go), [router_windows.go](../../internal/http/router_windows.go)). Il lit la configuration dans les variables globales du package `config` ([Chargement de la configuration](../configuration/chargement-et-reference.md)) et délègue la production du contenu à [heffalump](moteur-markov-heffalump.md).

## Enregistrement des routes

`Serve()` construit un routeur `fasthttp/router` selon deux modes exclusifs :

- **`catchall = true`** : une seule route `GET /{path:*}` capture tous les chemins. Aucun handler `robots.txt` n'est créé, même si `makerobots = true`.
- **`catchall = false`** : une route `GET /<chemin>` par entrée de `paths`, plus `GET /robots.txt` si `makerobots = true`. Tout autre chemin reçoit la réponse 404 du routeur.

Seules des routes `GET` existent. Par ailleurs le serveur est configuré avec `GetOnly: true`. Lors des essais manuels, un `POST` a reçu `400` et un `HEAD` a reçu `405`.

## `robots.txt`

Le handler génère le fichier à chaque requête : `User-agent: *` puis une ligne `Disallow: <chemin>` par entrée de `paths`, avec des fins de ligne `\r\n`. Il journalise `SERVE_ROBOTS` au niveau debug. L'intérêt est de rendre les chemins du piège « interdits » : seuls les robots qui ignorent le fichier y tombent. Le contenu est écrit avec `fmt.Fprint` et non `Fprintf`, de sorte qu'un `%` dans un chemin configuré ne corrompt pas la sortie.

## Traitement d'une requête piégée

Le handler `hellPot` enchaîne :

1. **Identité du client** : `REMOTE_ADDR` est la valeur de l'en-tête nommé par `http.real_ip_header` s'il est présent et non vide, sinon l'adresse TCP du pair. L'en-tête est cru tel quel : voir [CI, sécurité et releases](../operations/ci-securite-et-releases.md).
2. **Liste noire** : si le `User-Agent` contient (comparaison `strings.Contains`, sensible à la casse) l'une des chaînes de `uagent_string_blacklist`, la requête reçoit `404 Not found`, avec une ligne de log de niveau trace `Ignoring useragent`. Le client n'est pas piégé.
3. **Log `NEW`** (niveau info) avec `USERAGENT`, `REMOTE_ADDR`, `URL`. Si `trace` est activé, le champ `caller` (le chemin) est ajouté.
4. **Flux** : le corps de la réponse est un `SetBodyStreamWriter` qui rappelle `heffalump.WriteHell` jusqu'à la première erreur, additionne les octets écrits, puis journalise `FINISH` avec `BYTES` et `DURATION` (en millisecondes). Une erreur d'écriture est journalisée en trace sous `END_ON_ERR`. Voir [Journalisation](../operations/journalisation.md).

## Paramètres du serveur fasthttp

`getSrv` construit un `fasthttp.Server` avec : l'en-tête `Server` pris dans `deception.server_name`, un `ReadTimeout` de 5 s, un corps de requête limité à 1 Mio, `MaxRequestsPerConn` à 2, `DisableKeepalive`, `GetOnly`, `CloseOnShutdown`, `MaxConnsPerIP` pris dans `performance.max_conns_per_ip`, et `Concurrency` égale à `max_workers` seulement si `restrict_concurrency` est actif (sinon `fasthttp.DefaultConcurrency`). Le détail des effets sur la charge est dans [Limites de connexions et performance](../operations/limites-et-performance.md).

## Écoute : TCP ou socket Unix

- Par défaut (et toujours sous Windows), le serveur écoute en TCP sur `bind_addr:bind_port` via `srv.ListenAndServe` et journalise `Listening and serving HTTP...` avec l'adresse dans `caller`.
- Avec `use_unix_socket = true` (hors Windows), `unix_socket_path` ne doit pas être vide, sinon arrêt fatal. `listenOnUnixSocket` supprime d'abord un socket existant, crée le nouveau sous un `umask` restrictif (`0077`), puis applique `unix_socket_permissions` (lues en octal) avec `chmod`. Sous Windows la fonction existe mais renvoie une erreur.

**Mode socket Unix** : `listenOnUnixSocket` reçoit le `fasthttp.Server` construit par `getSrv` et appelle `srv.Serve(listener)`. Les mêmes réglages qu'en TCP s'appliquent donc : en-tête `Server` issu de `deception.server_name`, `ReadTimeout`, taille de corps, `GetOnly`, keep-alive désactivé. Deux particularités : `MaxConnsPerIP` ne s'applique pas, parce que fasthttp ne compte que les connexions qui ont une adresse TCP, et `REMOTE_ADDR` vaut `0.0.0.0` (adresse de pair absente) tant que le reverse proxy n'envoie pas l'en-tête `real_ip_header`. Constaté lors des essais : `Server: nginx` et un `POST` refusé en `400`, comme en TCP. Historique : dans le projet d'origine, ce mode utilisait un serveur fasthttp par défaut (`Server: fasthttp`, aucun réglage appliqué) ; le fork l'a corrigé et un test (`router_unix_test.go`) vérifie l'en-tête `Server` et le refus des méthodes autres que GET sur le socket.

## Arrêt

Une erreur de `http.Serve` est fatale (`log.Fatal`). Sur `SIGINT` ou `SIGTERM`, `main` journalise `Shutting down server...` et quitte : il n'attend pas la fin des connexions en cours. Voir [Vue d'ensemble de l'architecture](vue-densemble.md).
