---
type: operations
title: CI, sécurité et releases
description: Les trois workflows GitHub Actions du dépôt, les contrôles de sécurité automatiques, la mise à jour des dépendances par Dependabot, la procédure de publication par tag et la règle de nommage des versions.
tags: [ci, github-actions, securite, release, dependabot]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-79b37831c9c81206da1d88ec
    resource: repo://.github/dependabot.yml
  - id: openwiki-source-1bc83e59b3e3c64091f68169
    resource: repo://.github/workflows/docker.yml
  - id: openwiki-source-15f4e91280f1a3e266ffc6ab
    resource: repo://.github/workflows/go.yml
  - id: openwiki-source-6fc8389dad89bcb0319c264a
    resource: repo://.github/workflows/release-command.yml
  - id: openwiki-source-c34322bad84de7b6c0a25de1
    resource: repo://cmd/HellPot/HellPot.go
  - id: openwiki-source-e689a4a46f2ebef989178800
    resource: repo://internal/http/router.go
  - id: openwiki-source-012f2c78e3b1446dfc35803f
    resource: repo://Makefile
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# CI, sécurité et releases

Trois workflows vivent dans `.github/workflows/`, plus la configuration Dependabot. Le dépôt n'utilise pas de secret personnalisé : la publication de l'image emploie le `GITHUB_TOKEN` du workflow.

## Workflows

| Fichier | Nom | Déclenchement | Rôle |
| --- | --- | --- | --- |
| [go.yml](../../.github/workflows/go.yml) | Vibe Check | tout push, et pull request vers `main` | Contrôles qualité et sécurité |
| [docker.yml](../../.github/workflows/docker.yml) | Docker (GHCR) | push sur `main`, tag `v*`, pull request vers `main`, manuel, et chaque lundi à 06:00 UTC | Construit et publie l'image : voir [Image Docker](docker-et-publication.md) |
| [release-command.yml](../../.github/workflows/release-command.yml) | Build and Release | création d'une release GitHub, ou manuel | Binaires multi-plateformes |

## Contrôles du workflow « Vibe Check »

Dans l'ordre : `go vet`, `gosec`, `go test -race`, `go build`, `govulncheck`. La version de Go est lue dans `go.mod` (donc 1.26 ou plus récent). `gosec` et `govulncheck` sont installés à chaque exécution avec `@latest` : ils ne sont **pas épinglés**, ce qui garde leurs bases de règles à jour mais peut faire échouer le workflow sans qu'aucun code n'ait changé.

Le test unitaire exécuté, avec le détecteur de concurrence, est celui de [heffalump](../architecture/moteur-markov-heffalump.md). Le Dockerfile rejoue `go vet` et `go test` avant de compiler.

## Mise à jour des dépendances

Dependabot ouvre des pull requests groupées : modules Go et GitHub Actions chaque jour, images de base Docker chaque semaine. Le workflow Docker est lui aussi reconstruit chaque lundi, ce qui reprend les correctifs de la base distroless et de Go même quand aucun commit n'a eu lieu.

## Publier une version

1. Choisir un tag qui commence par **un caractère de préfixe**, en pratique `v1.2.3`. Le programme retire le premier caractère de la version fournie par le linker : un tag `1.2.3` s'afficherait `.2.3`.
2. Pousser le tag : le workflow Docker publie les tags d'image versionnés.
3. Créer la release GitHub à partir de ce tag : le workflow « Build and Release » compile et attache les binaires, avec leur somme SHA-256.

Matrice des binaires : Linux, Windows, macOS (darwin) et FreeBSD en `386`, `amd64`, `arm64`, sauf darwin/`386` et windows/`arm64`. Les binaires sont compilés sans CGO et reçoivent le nom du tag dans `main.version`. Dans le projet d'origine, le workflow y plaçait la référence git complète (`refs/tags/...`), ce qui donnait une version erronée à l'affichage ; le fork utilise le nom du tag.

La version affichée par `--banner` est donc le tag sans son premier caractère ; sans tag, c'est le hash court du commit ou `dev`. Le Makefile suit la même logique : `make build` prend le dernier tag git trié par version. Voir aussi [Vue d'ensemble](../architecture/vue-densemble.md).

## Limites de confiance

- **`X-Real-IP` est cru tel quel.** HellPot lit l'IP du client dans l'en-tête configuré par `real_ip_header` sans en vérifier l'origine. Exposé directement sur Internet, un client peut donc falsifier le `REMOTE_ADDR` journalisé. Le service doit être derrière un reverse proxy qui écrase cet en-tête, ou sur un réseau interne. Voir [Serveur HTTP](../architecture/serveur-http-et-routage.md).
- L'image s'exécute en root mais n'a besoin d'aucune capability ni d'un système de fichiers inscriptible hors `/logs`.
- Le dépôt n'a pas de fichier `SECURITY.md` : les vulnérabilités connues des dépendances sont détectées automatiquement par `govulncheck`.
