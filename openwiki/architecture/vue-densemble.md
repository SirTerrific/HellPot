---
type: architecture
title: Vue d'ensemble de l'architecture
description: Les cinq packages de HellPot, leur rôle, et la séquence de démarrage du processus, de l'initialisation de la configuration et du logger jusqu'au serveur HTTP et à l'arrêt sur signal.
tags: [architecture, demarrage, packages, cycle-de-vie]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-c34322bad84de7b6c0a25de1
    resource: repo://cmd/HellPot/HellPot.go
  - id: openwiki-source-a8910515ddd14810ad43f5c1
    resource: repo://internal/config/config.go
  - id: openwiki-source-f715782272efaf73ca0c929d
    resource: repo://internal/config/globals.go
  - id: openwiki-source-091ee1219a9d0c94513b0e00
    resource: repo://internal/config/logger.go
  - id: openwiki-source-f4868503981e3b707b159bb2
    resource: repo://internal/extra/banner.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Vue d'ensemble de l'architecture

HellPot est un binaire Go unique (module `github.com/SirTerrific/HellPot`) qui répond à des requêtes HTTP par un flux de texte infini. Il n'a ni base de données ni état persistant autre que ses fichiers de log.

## Packages et responsabilités

| Package | Rôle | Page détaillée |
| --- | --- | --- |
| `cmd/HellPot` | Point d'entrée : initialisation (config, logger, bannière) puis lancement du serveur et attente d'un signal | cette page |
| `internal/config` | Lecture de la configuration (fichier TOML, variables d'environnement, options de ligne de commande), variables globales exportées, démarrage du logger | [Configuration](../configuration/chargement-et-reference.md), [Journalisation](../operations/journalisation.md) |
| `internal/http` | Serveur fasthttp, routes, `robots.txt`, handler du piège, écoute TCP ou socket Unix | [Serveur HTTP](serveur-http-et-routage.md) |
| `heffalump` | Génération du texte infini par chaîne de Markov | [Moteur de Markov](moteur-markov-heffalump.md) |
| `internal/extra` | Affichage de la bannière | cette page |

Le sens des dépendances est `cmd` vers `config`, `extra` et `http` ; `http` vers `heffalump` et `config` ; `heffalump` vers `config` (seulement pour le logger). Les packages ne se passent pas la configuration en paramètre : ils lisent des variables globales de `config`, remplies une fois pour toutes pendant `config.Init()`.

## Séquence de démarrage

Tout le démarrage a lieu dans la fonction `init()` de `cmd/HellPot`, avant `main()` :

1. **Version** : si le linker a fourni `main.version`, `config.Version` reçoit cette valeur **sans son premier caractère** (le `v` d'un tag comme `v1.2.3`). Sinon `config` utilise le hash court (7 caractères) du commit lu dans les informations de build, ou `dev`.
2. **`config.Init()`** : analyse des options de ligne de commande, chargement du fichier, surcharge par variables d'environnement. Détails dans [Chargement de la configuration](../configuration/chargement-et-reference.md).
3. **`--banner`** : si demandé, affiche la bannière et quitte avec le code 0, sans démarrer de logger.
4. **Logger** : avec `docker_logging`, les logs JSON vont uniquement sur la sortie standard (couleurs coupées) ; sinon un fichier de log est ouvert et la console reçoit en plus une version lisible. Voir [Journalisation](../operations/journalisation.md).
5. **Bannière** : sous Windows ou avec `nocolor`, une simple ligne `HellPot <version>` ; sinon une bannière ANSI décompressée à partir d'une chaîne embarquée, dont les couleurs sont tirées au hasard avec `crypto/rand`.
6. **Lignes de démarrage** : le fichier de configuration (`caller: config`), le fichier de log (`caller: logger`), puis les messages `debug enabled` et `trace enabled`.

Le package `heffalump` a sa propre initialisation, indépendante : sa table de Markov est construite quand le package est chargé, donc avant même `config.Init()`. Une défaillance à ce stade fait paniquer le processus.

## Exécution et arrêt

`main()` lance `http.Serve()` dans une goroutine : si cette fonction retourne une erreur (adresse déjà utilisée, par exemple), le processus s'arrête avec un `log.Fatal`. Le thread principal attend `SIGINT` ou `SIGTERM`, journalise `Shutting down server...` et retourne. Il n'y a pas d'arrêt gracieux : les connexions en cours sont coupées. C'est cohérent avec le rôle du service, dont les connexions sont conçues pour ne jamais finir.

## Le logger partagé

Les packages `http` et `heffalump` obtiennent un pointeur vers le logger global avec `config.GetLogger()`. C'est un pointeur vers une variable de package : le logger est affecté plus tard par `StartLogger`, et tous les pointeurs déjà pris voient la valeur à jour. Tant que `StartLogger` n'a pas été appelé, ce logger est la valeur zéro de zerolog.
