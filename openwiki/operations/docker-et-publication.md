---
type: operations
title: Image Docker et publication sur GHCR
description: Comment l'image HellPot est construite (Dockerfile multi-étapes, compilation croisée amd64 et arm64), ce qu'elle contient, comment elle est publiée sur GHCR et comment la déployer avec docker run ou docker compose.
tags: [docker, ghcr, distroless, multi-arch, deploiement]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-1bc83e59b3e3c64091f68169
    resource: repo://.github/workflows/docker.yml
  - id: openwiki-source-bb1ebe868e35e9e500714501
    resource: repo://Dockerfile
  - id: openwiki-source-499d45fbdac8421cc7917be2
    resource: repo://internal/config/arguments.go
  - id: openwiki-source-091ee1219a9d0c94513b0e00
    resource: repo://internal/config/logger.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Image Docker et publication sur GHCR

L'image est publiée sur `ghcr.io/sirterrific/hellpot` pour `linux/amd64` et `linux/arm64`. Sa disposition est volontairement identique à celle de l'image du projet d'origine : voir [Fork et compatibilité](../guides/fork-et-compatibilite.md).

## Dockerfile

Deux étapes ([Dockerfile](../../Dockerfile)).

**Étape de build** (`golang:1.27`, exécutée sur la plateforme de la machine qui construit, `$BUILDPLATFORM`) :

1. `go mod download` sur les seuls `go.*`, pour profiter du cache de couches.
2. Copie de tout le contexte, puis **`go vet` et `go test`** : un test qui échoue fait échouer la construction de l'image.
3. Compilation avec `CGO_ENABLED=0`, `-trimpath` et `-s -w`, avec `GOOS` et `GOARCH` pris dans les arguments `TARGETOS` et `TARGETARCH` de buildx. La compilation croisée est donc native : l'arm64 est produit sans émulation QEMU.
4. La version injectée (`-X main.version`) est l'argument de build `VERSION` ; à défaut, le dernier tag git du contexte. Le dépôt n'a pas de `.dockerignore`, donc `.git` fait partie du contexte, ce qui permet à `git tag` de fonctionner, et à Go d'enregistrer le commit pour les builds sans tag. Voir [CI, sécurité et releases](ci-securite-et-releases.md).

**Étape finale** (`gcr.io/distroless/static-debian13`) : le binaire est copié en `/app`, [docker_config.toml](../../docker_config.toml) en `/config`, le port 8080 est déclaré et l'entrypoint est `/app -c /config`. Le label `org.opencontainers.image.source` pointe vers ce dépôt. Aucune directive `USER` : le conteneur s'exécute en root. L'image n'a pas de shell, donc ni `docker exec ... sh` ni healthcheck basé sur des commandes shell.

## Configuration embarquée

`/config` est chargé avec `-c`, donc **sans valeurs par défaut intégrées** : voir [Chargement de la configuration](../configuration/chargement-et-reference.md). Pour utiliser un autre fichier, le monter sur `/config`. Si on ajoute `-c autre.toml` aux arguments, il est fusionné par-dessus `/config` car l'entrypoint le charge déjà. Les variables `HELLPOT_*` surchargent le tout. `docker_config.toml` active `catchall`, écoute `0.0.0.0:8080` et écrit les logs dans `/logs/`.

## Publication : workflow `docker.yml`

[docker.yml](../../.github/workflows/docker.yml) construit pour `linux/amd64` et `linux/arm64` avec buildx, passe `VERSION` seulement quand l'événement est un tag, et utilise le cache GitHub Actions. Il se connecte à GHCR avec le `GITHUB_TOKEN` (permission `packages: write`) et ne **pousse pas** pour les pull requests : elles construisent seulement.

| Tag d'image | Condition |
| --- | --- |
| `latest` | exécution sur la branche par défaut (push sur `main`, et reconstruction du lundi) |
| `sha-<commit court>` | tout événement qui pousse |
| `<version>` et `<majeure>.<mineure>` | push d'un tag semver `v*`, par exemple `v1.2.3` donne `1.2.3` et `1.2` |

La reconstruction hebdomadaire reprend les correctifs de sécurité de la base distroless et de Go sans changement de code.

## Déploiement

```
docker run -d --name hellpot -p 8080:8080 -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest
```

Derrière un reverse proxy, on utilise `expose` plutôt que `ports` et on relève `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP` (tous les clients partagent l'IP du proxy) : voir [Limites de connexions et performance](limites-et-performance.md).

Le durcissement `read_only: true`, `cap_drop: [ALL]` et `security_opt: [no-new-privileges:true]` est vérifié avec cette image. Avec `read_only`, un volume monté sur `/logs` est indispensable, sauf si `docker_logging` est activé (logs JSON sur la sortie standard, sans fichier) : le logger crée son fichier au démarrage et le processus s'arrête avec « cannot create log file » s'il n'y parvient pas. La variable `TZ` règle le fuseau du nom de fichier de log et du champ `time`. Les logs eux-mêmes sont décrits dans [Journalisation](journalisation.md).
