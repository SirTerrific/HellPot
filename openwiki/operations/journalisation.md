---
type: operations
title: Journalisation et format des logs
description: Comment HellPot écrit ses logs (JSON zerolog dans un fichier ou sur stdout, console lisible), le nom du fichier, le catalogue des messages et de leurs champs, et la contrainte de ne pas changer ce format.
tags: [logs, zerolog, json, observabilite, compatibilite]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-c34322bad84de7b6c0a25de1
    resource: repo://cmd/HellPot/HellPot.go
  - id: openwiki-source-a8910515ddd14810ad43f5c1
    resource: repo://internal/config/config.go
  - id: openwiki-source-091ee1219a9d0c94513b0e00
    resource: repo://internal/config/logger.go
  - id: openwiki-source-a9ca8194716b05765ae68c77
    resource: repo://internal/http/robots.go
  - id: openwiki-source-e689a4a46f2ebef989178800
    resource: repo://internal/http/router.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Journalisation et format des logs

Les logs sont la sortie principale de HellPot : chaque bot piégé y laisse une trace, souvent exploitée par un collecteur externe (Splunk, SCOM...). **Le format est donc un contrat de compatibilité** : toute modification des clés, des niveaux, du texte des messages ou de l'ordre des champs doit être évitée. Voir [Fork et compatibilité](../guides/fork-et-compatibilite.md).

## Initialisation

`StartLogger` ([logger.go](../../internal/config/logger.go)) est appelé une fois, au démarrage, par `cmd/HellPot` (voir [Vue d'ensemble](../architecture/vue-densemble.md)). Il construit un logger zerolog avec horodatage et le range dans une variable globale que tous les packages lisent avec `GetLogger()`.

Deux modes :

- **Mode normal** : un fichier de log JSON est ouvert (création si besoin, ajout à la fin, droits `0666` sous réserve de l'umask), **et** la console reçoit une version lisible (`ConsoleWriter`, format d'heure `console_time_format`, couleurs sauf `nocolor`).
- **Mode `docker_logging`** : le JSON est écrit uniquement sur la sortie standard. Aucun fichier, aucune couleur. C'est le mode pensé pour `docker logs` et les collecteurs qui lisent stdout. Dans ce mode, la ligne de démarrage qui nomme le fichier de log affiche `/dev/stdout`.

## Fichier de log

- **Répertoire** : `logger.directory`, ou à défaut `$HOME/.local/share/HellPot/logs`. Il est créé avec les droits `0750`. L'image Docker utilise `/logs/`.
- **Nom** : `HellPot.log`, ou avec `use_date_filename` `HellPot_<date>.log`, où la date est le format RFC822 de l'heure de **démarrage** dont les espaces deviennent `_` et les `:` deviennent `-`, par exemple `HellPot_03_Oct_26_20-56_EDT.log`. Le nom est fixé une fois : un processus de longue durée écrit toujours dans le même fichier, et chaque redémarrage en crée un nouveau.
- **Fuseau** : celui du processus, donc la variable `TZ` en conteneur. Il apparaît dans le nom du fichier et dans le champ `time` de chaque ligne.
- Si le fichier ne peut pas être créé, le processus s'arrête immédiatement (« cannot create log file »).

## Catalogue des messages

Une ligne JSON a toujours `level`, ses champs propres, `time`, puis `message`. Les champs de requête sont en majuscules.

| Message | Niveau | Émetteur | Champs |
| --- | --- | --- | --- |
| nom du fichier de configuration | info | démarrage | `caller: "config"`, `file` |
| chemin du fichier de log | info | démarrage | `caller: "logger"` |
| `debug enabled` | debug | démarrage | `caller: "logger"` |
| `trace enabled` | trace | démarrage | `caller: "logger"` |
| `Add route: <chemin>` | trace | routeur | `caller: "router"` |
| `Catch-All mode enabled...` | trace | routeur | |
| `Listening and serving HTTP...` | info | routeur | `caller` : adresse TCP ou chemin du socket |
| `unix_socket_path configuration directive appears to be empty` | fatal | routeur | |
| `NEW` | info | requête piégée | `USERAGENT`, `REMOTE_ADDR`, `URL` |
| `FINISH` | info | fin de requête | mêmes champs, plus `BYTES` et `DURATION` (millisecondes) |
| `END_ON_ERR` | trace | fin de requête | mêmes champs, plus `error` |
| `Ignoring useragent` | trace | liste noire | `USERAGENT`, `REMOTE_ADDR`, `URL` |
| `SERVE_ROBOTS` | debug | `robots.txt` | mêmes champs, plus `PATHS` |
| `SERVE_ROBOTS_ERROR` | error | `robots.txt` | mêmes champs, plus `error` |
| `panic recovered!` | error | heffalump | `caller` : la valeur de la panique |
| `Shutting down server...` | warn | arrêt | |
| `HTTP error` | fatal | arrêt | `error` |

`caller` s'ajoute aussi aux lignes de requête (valeur : le chemin) quand `trace` est activé.

Les messages internes de fasthttp passent par le même logger (`Logger: log` dans la configuration du serveur) et sortent au niveau **debug**, avec le texte de la bibliothèque comme `message` : par exemple `error when serving connection ...` ou `The number of connections from <ip> exceeds MaxConnsPerIP=<n>`, ce dernier signalant un client refusé en HTTP 429. Voir [Limites de connexions et performance](limites-et-performance.md). Le texte de ces messages vient de fasthttp et peut changer d'une version à l'autre de la bibliothèque.

## Niveaux : debug et trace ne filtrent rien

Les options `debug` et `trace` ne font que **relever** le niveau global de zerolog ; sans elles, le niveau global n'est pas touché et les lignes debug et trace sont écrites quand même. C'est pourquoi une configuration avec `debug = false` et `trace = false` produit tout de même des lignes `TRC` et `DBG`. Ce comportement vient du projet d'origine et est conservé pour que le contenu des logs reste identique. Voir [Chargement de la configuration](../configuration/chargement-et-reference.md).

## Exemple

```json
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","time":"2026-10-03T22:00:36Z","message":"NEW"}
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","BYTES":179549,"DURATION":15.914698,"time":"2026-10-03T22:00:36Z","message":"FINISH"}
```

`REMOTE_ADDR` est la valeur de l'en-tête `real_ip_header` si elle est présente, sinon l'adresse du pair : voir [Serveur HTTP](../architecture/serveur-http-et-routage.md).
