---
type: configuration
title: Chargement de la configuration, options et variables d'environnement
description: D'où HellPot lit sa configuration (fichier TOML, variables HELLPOT_, options de ligne de commande), dans quel ordre, pourquoi les valeurs par défaut disparaissent avec -c, et les particularités connues du code.
tags: [configuration, toml, variables-environnement, cli, koanf]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-c239e6e06e2b1ec445a5f81f
    resource: repo://docker_config.toml
  - id: openwiki-source-499d45fbdac8421cc7917be2
    resource: repo://internal/config/arguments.go
  - id: openwiki-source-a8910515ddd14810ad43f5c1
    resource: repo://internal/config/config.go
  - id: openwiki-source-871a8611f939f64b133e20a8
    resource: repo://internal/config/defaults.go
  - id: openwiki-source-f715782272efaf73ca0c929d
    resource: repo://internal/config/globals.go
  - id: openwiki-source-cd862ac729c594a324d155a6
    resource: repo://internal/config/help.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Chargement de la configuration, options et variables d'environnement

Tout le code est dans `internal/config` ([config.go](../../internal/config/config.go), [defaults.go](../../internal/config/defaults.go), [arguments.go](../../internal/config/arguments.go), [globals.go](../../internal/config/globals.go), [help.go](../../internal/config/help.go)). La bibliothèque est koanf (importée sous le nom `viper`, instance `snek`). Le résultat est un ensemble de **variables globales** (`HTTPBind`, `Paths`, `MaxConnsPerIP`...) que les autres packages lisent directement : voir [Vue d'ensemble](../architecture/vue-densemble.md).

## Couches de configuration

Dans l'ordre d'application, la dernière couche gagne :

1. **Valeurs par défaut intégrées** (`defOpts`), **uniquement si `-c` n'est pas utilisé**.
2. **Fichier TOML**.
3. **Variables d'environnement `HELLPOT_*`**, dans les deux modes.
4. **Options `-v`, `-vv`, `--nocolor`**, qui ne peuvent qu'activer quelque chose.

## Sans `-c` : recherche du fichier

`Init()` applique d'abord les valeurs par défaut, puis choisit un fichier :

1. `/etc/HellPot/config.toml`, s'il existe (jamais sous Windows) ;
2. sinon `config.toml` dans le répertoire de configuration utilisateur (`os.UserConfigDir()`, sous-dossier `HellPot`, créé au besoin) ;
3. repli sur `./config.toml` seulement si aucun répertoire utilisateur n'est déterminable.

Si le fichier choisi est absent ou illisible, HellPot affiche « No configuration file found, writing new configuration file... », écrit un nouveau fichier (droits `0600`) à l'emplacement du répertoire utilisateur à partir des valeurs par défaut, puis le charge. Une erreur d'écriture ou de chargement termine le processus avec le code 1.

## Avec `-c` : pas de valeurs par défaut

`-c <fichier>` (ou `--config`) charge le fichier immédiatement, puis `Init()` **retourne sans appeler `setDefaults()`**. Conséquences :

- Toute clé absente du fichier prend sa valeur zéro (chaîne vide, `false`, `0`, liste vide). Un fichier qui ne définit que `[http]` n'enregistre aucune route : tous les chemins répondent 404 et il n'y a pas de `robots.txt`. De même `use_date_filename` absent donne un fichier de log sans date.
- **Exception** : `performance.max_conns_per_ip` garde sa valeur `10` quand la clé est absente, car la variable `MaxConnsPerIP` est initialisée à 10 et n'est écrasée que si la clé existe.
- `--genconfig` est silencieusement ignoré : il n'agit que dans `setDefaults()`.
- Un `-c` sans fichier après lui est une erreur fatale ; un fichier illisible aussi.
- `-c` peut être répété : chaque fichier est fusionné par-dessus le précédent. C'est ce qui se produit dans l'image Docker quand on ajoute `-c` aux arguments : l'entrypoint charge déjà `/config`.

## Variables d'environnement

`associateExportedVariables` charge les variables préfixées `HELLPOT_` **après** le fichier, dans les deux modes. Transformation du nom : retrait du préfixe, minuscules, puis `__` devient un caractère réservé, chaque `_` restant devient `.`, et le caractère réservé redevient `_`. Autrement dit : un `_` sépare les niveaux, `__` représente un `_` du nom de la clé.

| Clé | Variable |
| --- | --- |
| `http.bind_addr` | `HELLPOT_HTTP_BIND__ADDR` |
| `http.router.catchall` | `HELLPOT_HTTP_ROUTER_CATCHALL` |
| `logger.docker_logging` | `HELLPOT_LOGGER_DOCKER__LOGGING` |
| `performance.max_conns_per_ip` | `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP` |

## Correspondance clé vers variable globale

`processOpts` copie chaque clé dans une variable : `http.bind_addr` (`HTTPBind`), `http.bind_port` (`HTTPPort`), `http.real_ip_header` (`HeaderName`), `logger.directory` (`logDir`), `logger.console_time_format` (`ConsoleTimeFormat`), `deception.server_name` (`FakeServerName`), `http.router.paths` (`Paths`), `http.uagent_string_blacklist` (`UseragentBlacklistMatchers`), `http.use_unix_socket`, `logger.debug`, `logger.trace`, `logger.nocolor`, `logger.docker_logging`, `http.router.makerobots`, `http.router.catchall`, `performance.restrict_concurrency`, `performance.max_workers`. Cas particuliers traités après :

- `max_conns_per_ip` : voir ci-dessus.
- Socket Unix : `unix_socket_path` et `unix_socket_permissions` ne sont lus que si `use_unix_socket` est vrai. Les permissions sont interprétées en **octal** ; une valeur invalide est ignorée sans message et laisse la permission à 0, ce qui rendrait le socket inutilisable.

## Valeurs par défaut

Sans `-c` : debug actif, trace inactif, `use_date_filename` actif, `docker_logging` inactif, `console_time_format` `3:04PM`, écoute `127.0.0.1:8080`, `real_ip_header` `X-Real-IP`, socket Unix désactivé (`/var/run/hellpot`, permissions `0666`), `catchall` inactif, `makerobots` actif, `paths` `wp-login.php` et `wp-login`, liste noire `Cloudflare-Traffic-Manager`, `restrict_concurrency` inactif, `max_workers` 256, `max_conns_per_ip` 10, `server_name` `nginx`.

## Configuration de l'image Docker

[docker_config.toml](../../docker_config.toml) est chargé avec `-c`, donc **sans** valeurs par défaut. Il fixe : `server_name` `nginx`, écoute `0.0.0.0:8080`, `real_ip_header` `X-Real-IP`, liste noire `Cloudflare-Traffic-Manager` et `curl`, `catchall` actif, `debug` et `trace` à `false`, répertoire de logs `/logs/`, `nocolor` actif, `use_date_filename` actif. Tout le reste (par exemple `makerobots`, `docker_logging`) vaut zéro. Voir [Image Docker](../operations/docker-et-publication.md).

## Options de ligne de commande

`argParse` parcourt tous les arguments : `--debug` ou `-v`, `--trace` ou `-vv`, `--nocolor`, `--banner`, `--genconfig` activent un booléen ; `-c` ou `--config` charge un fichier ; `-h` affiche l'aide.

- Le texte d'aide liste `--help`, mais **seul `-h` est reconnu** : `--help` est ignoré et le serveur démarre.
- L'aide ne s'affiche que si la sortie standard est un terminal ; sinon le processus sort avec le code 1.
- L'aide annonce que `--genconfig` écrit `HellPot.toml` ; en réalité le fichier est `./config.toml`, avec des droits `0600`, puis le processus quitte avec le code 0.
- `-v` et `-vv` (ou les clés `debug` et `trace`) règlent le niveau global de zerolog sur debug ou trace, et ne le baissent jamais. Sans eux, le niveau global n'est pas modifié : c'est pourquoi les lignes debug et trace apparaissent même avec `debug = false`. Voir [Journalisation](../operations/journalisation.md).
