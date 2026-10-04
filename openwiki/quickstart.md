---
type: "Référence"
title: "Démarrage rapide et carte de la documentation"
openwiki_generated: true
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-c34322bad84de7b6c0a25de1
    resource: repo://cmd/HellPot/HellPot.go
  - id: openwiki-source-c239e6e06e2b1ec445a5f81f
    resource: repo://docker_config.toml
  - id: openwiki-source-bb1ebe868e35e9e500714501
    resource: repo://Dockerfile
  - id: openwiki-source-a8910515ddd14810ad43f5c1
    resource: repo://internal/config/config.go
  - id: openwiki-source-871a8611f939f64b133e20a8
    resource: repo://internal/config/defaults.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---


# Démarrage rapide et carte de la documentation

HellPot est un honeypot HTTP « sans fin » : un client qui visite un chemin piégé reçoit un flux infini de texte pseudo-aléatoire (chaîne de Markov sur un texte de Nietzsche) et reste bloqué tant qu'il lit. Ce dépôt est un fork maintenu de `yunginnanet/HellPot`, dont le projet d'origine n'est plus maintenu. Il vise un remplacement direct : même configuration, mêmes logs. Voir [Fork et compatibilité](guides/fork-et-compatibilite.md).

## Le lancer

- **Docker** (recommandé) : `docker pull ghcr.io/sirterrific/hellpot:latest`, puis `docker run -d -p 8080:8080 -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest`. L'image écoute sur le port 8080, lit sa configuration dans `/config` et écrit ses logs dans `/logs`. Détails : [Image Docker](operations/docker-et-publication.md).
- **Binaire** : télécharger une release GitHub, ou compiler avec Go 1.26 ou plus récent (`make`). Sans option, il crée sa configuration par défaut dans le répertoire de configuration de l'utilisateur et écoute sur `127.0.0.1:8080`.
- **Fichier de configuration** : `./HellPot --genconfig` écrit `config.toml` ; `./HellPot -c config.toml` le charge. Attention : avec `-c`, les valeurs par défaut ne s'appliquent plus. Voir [Chargement de la configuration](configuration/chargement-et-reference.md).

Le point d'entrée du programme est [HellPot.go](../cmd/HellPot/HellPot.go) : il initialise la configuration, le logger et la bannière, puis lance le serveur HTTP et attend un signal d'arrêt.

## Quelle page lire

| Je veux... | Page |
| --- | --- |
| comprendre l'organisation du code et le démarrage | [Vue d'ensemble](architecture/vue-densemble.md) |
| savoir quelles requêtes sont piégées, comment marchent `robots.txt`, la liste noire, le socket Unix | [Serveur HTTP, routage et piège](architecture/serveur-http-et-routage.md) |
| comprendre comment le texte infini est généré | [Moteur de Markov](architecture/moteur-markov-heffalump.md) |
| régler une option, une variable d'environnement, comprendre `-c` | [Chargement de la configuration](configuration/chargement-et-reference.md) |
| lire ou parser les logs, comprendre leur format | [Journalisation](operations/journalisation.md) |
| comprendre des refus HTTP 429 ou dimensionner le service | [Limites de connexions et performance](operations/limites-et-performance.md) |
| construire, publier ou déployer l'image | [Image Docker et publication](operations/docker-et-publication.md) |
| publier une version, comprendre la CI et la sécurité | [CI, sécurité et releases](operations/ci-securite-et-releases.md) |
| savoir ce qui a changé par rapport à l'upstream et ce qui ne doit pas changer | [Fork et compatibilité](guides/fork-et-compatibilite.md) |

## Pièges fréquents

- Derrière un reverse proxy, la limite par défaut de 10 connexions par IP plafonne toute l'instance : relever `performance.max_conns_per_ip`.
- `debug = false` ne supprime pas les lignes debug et trace des logs.
- `--help` n'est pas reconnu, seul `-h` l'est.
- `X-Real-IP` est cru tel quel : ne pas exposer HellPot directement sur Internet si l'IP journalisée doit être fiable.
- Les fichiers `README.md` (anglais) et `README.fr.md` (français) à la racine documentent l'usage ; [CHANGELOG.md](../CHANGELOG.md) liste les changements du fork.
