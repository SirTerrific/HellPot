---
type: guide
title: Fork et garanties de compatibilité
description: Ce qui distingue ce fork du projet yunginnanet/HellPot, ce qui doit rester identique pour ne rien casser (config, logs, image) et comment passer d'une image construite localement à l'image publiée sur GHCR.
tags: [fork, compatibilite, migration, docker, logs]
sources:
  - id: openwiki-source-ca6cb4b1a14fd7969dfae3ec
    resource: repo://CHANGELOG.md
  - id: openwiki-source-bb1ebe868e35e9e500714501
    resource: repo://Dockerfile
  - id: openwiki-source-7bd911fdd3026b7b031a01e3
    resource: repo://go.mod
  - id: openwiki-source-871a8611f939f64b133e20a8
    resource: repo://internal/config/defaults.go
  - id: openwiki-source-f715782272efaf73ca0c929d
    resource: repo://internal/config/globals.go
  - id: openwiki-source-fd1f4e266c537f5376b692a9
    resource: repo://internal/http/router_unix.go
generated: { by: "claude-code", at: "2026-10-04T19:20:53.104Z" }
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T19:20:53.104Z
---

# Fork et garanties de compatibilité

Ce dépôt est un fork maintenu de [yunginnanet/HellPot](https://github.com/yunginnanet/HellPot), dont le projet d'origine n'est plus maintenu. Il part du commit upstream `0ba62c9`. Aucune version n'a encore été étiquetée : tout est listé sous « Unreleased » dans [CHANGELOG.md](../../CHANGELOG.md).

## Règle directrice : remplacement direct

Le fork vise un remplacement direct. Ces éléments doivent rester identiques, et tout changement futur doit être mesuré contre eux :

- **Clés de configuration et valeurs par défaut** : `defOpts` dans [defaults.go](../../internal/config/defaults.go). Une nouvelle clé doit avoir une valeur par défaut qui reproduit l'ancien comportement. Exemple : `performance.max_conns_per_ip` vaut 10, la valeur qui était codée en dur. Voir [Chargement de la configuration](../configuration/chargement-et-reference.md).
- **Structure des logs** : clés JSON, niveaux, messages et ordre des champs. Catalogue dans [Journalisation](../operations/journalisation.md).
- **Comportement HTTP** : routes, `robots.txt`, liste noire, en-tête `Server` (à l'exception du mode socket Unix, corrigé : voir ci-dessous). Voir [Serveur HTTP](../architecture/serveur-http-et-routage.md).
- **Disposition de l'image** : binaire `/app`, configuration `/config`, logs `/logs`, port 8080, entrypoint `/app -c /config`. Voir [Image Docker](../operations/docker-et-publication.md).

## Ce qui a changé par rapport à l'upstream

- Dépendances Go mises à jour (fasthttp, zerolog, koanf, `x/sys`, `x/term`) ; la directive `go` passe à 1.26.0, exigée par ces dépendances.
- Module renommé en `github.com/SirTerrific/HellPot` : le chemin d'installation `go install` change, le binaire non.
- Dockerfile : image de build Go 1.27, image d'exécution `distroless/static-debian13`, compilation croisée native pour amd64 et arm64, argument de build `VERSION`.
- Publication de l'image sur `ghcr.io/sirterrific/hellpot` par GitHub Actions ; suppression des workflows qui visaient le dépôt et le Docker Hub de l'upstream.
- Nouveau paramètre optionnel `performance.max_conns_per_ip`.
- Corrections : le mode socket Unix utilise désormais le serveur configuré (en-tête `Server`, délais, GET seulement) au lieu d'un serveur fasthttp par défaut — seul changement de comportement volontaire, qui ne concerne que les utilisateurs de `use_unix_socket` ; `robots.txt` n'utilise plus une chaîne de format non constante ; le pool de buffers inutilisé de `heffalump` est supprimé.
- Premier test unitaire, CI renforcée (`govulncheck`), documentation en anglais et en français.

## Différences de texte dans les logs

Deux messages de niveau debug ou trace changent, parce que leur texte vient de la bibliothèque fasthttp mise à jour et non de HellPot : `END_ON_ERR` rapporte `fasthttputil: connection closed` (avant : `connection closed`) et la ligne de requête non-GET devient `fasthttp: non-get request received`. Tous les champs et messages propres à HellPot sont inchangés. Un filtre Splunk ou SCOM basé sur le texte de ces deux messages devrait être revu.

## Particularités de l'upstream conservées volontairement

Ces comportements existaient déjà et n'ont pas été modifiés, pour ne pas changer la sortie ni la configuration existantes : les lignes debug et trace écrites même avec `debug = false`, `--help` non reconnu (seul `-h` l'est), l'aide qui annonce `HellPot.toml` alors que `--genconfig` écrit `config.toml`, et l'absence de valeurs par défaut avec `-c`. Chacun est expliqué sur la page de son sujet.

## Migrer d'une image construite localement vers GHCR

Cas d'usage : un `docker-compose.yml` avec un bloc `build:` pointant vers un clone local et `image: hellpot:local`.

1. Remplacer le bloc `build:` et la ligne `image:` par `image: ghcr.io/sirterrific/hellpot:latest`.
2. Conserver tels quels les volumes (`./logs:/logs`), `TZ`, le réseau et le reste.
3. Lancer `docker compose pull`, puis `docker compose up -d`.

Rien d'autre ne change : même chemin de logs, même format, même fichier de configuration embarqué. Avec `TZ=America/Toronto`, l'ancienne et la nouvelle image produisent le même nom de fichier de log (par exemple `HellPot_03_Oct_26_20-56_EDT.log`) et le même fuseau dans le champ `time`. Si HellPot tourne derrière un reverse proxy, penser à relever `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP` : voir [Limites de connexions et performance](../operations/limites-et-performance.md).

## Version affichée

Le Dockerfile prend `VERSION` en argument de build et, à défaut, le dernier tag git du contexte. Tant qu'aucun tag n'existe, la version affichée par `--banner` est le hash court du commit ou `dev`. Pousser un tag `v*` fait publier des tags d'image versionnés. Voir [CI, sécurité et releases](../operations/ci-securite-et-releases.md).
