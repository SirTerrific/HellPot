---
type: operations
title: Limites de connexions et performance
description: Les réglages qui bornent la charge de HellPot (connexions par IP, concurrence, timeouts), l'effet d'un reverse proxy sur la limite par IP, ce qui a été mesuré, et pourquoi le pool de buffers a été retiré.
tags: [performance, limites, maxconnsperip, concurrence, reverse-proxy]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T02:02:24.687Z
sources:
  - id: openwiki-source-6225717d77c0df8b3a43d4df
    resource: repo://heffalump/heffalump.go
  - id: openwiki-source-a8910515ddd14810ad43f5c1
    resource: repo://internal/config/config.go
  - id: openwiki-source-871a8611f939f64b133e20a8
    resource: repo://internal/config/defaults.go
  - id: openwiki-source-f715782272efaf73ca0c929d
    resource: repo://internal/config/globals.go
  - id: openwiki-source-fd1f4e266c537f5376b692a9
    resource: repo://internal/http/router_unix.go
  - id: openwiki-source-e689a4a46f2ebef989178800
    resource: repo://internal/http/router.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Limites de connexions et performance

HellPot est un tarpit : chaque bot piégé garde une **connexion ouverte indéfiniment**, car la réponse ne se termine que lorsque le client part (voir [Moteur de Markov](../architecture/moteur-markov-heffalump.md)). Le nombre de connexions simultanées est donc le nombre de bots actuellement piégés, et les limites ci-dessous décident combien le service en accepte.

## Réglages

Tous sont posés par `getSrv` dans [router.go](../../internal/http/router.go) pour le mode TCP.

| Réglage | Valeur | Effet |
| --- | --- | --- |
| `performance.max_conns_per_ip` | 10 par défaut, `0` = illimité | Connexions simultanées maximales par IP distante ; au-delà, le client reçoit HTTP 429 |
| `performance.restrict_concurrency` et `max_workers` | inactif ; 256 | `max_workers` ne s'applique que si `restrict_concurrency` est actif ; sinon la concurrence est `fasthttp.DefaultConcurrency` |
| `ReadTimeout` | 5 s | Borne la lecture de la requête (pas l'écriture de la réponse) |
| `MaxRequestBodySize` | 1 Mio | Plafond du corps de requête |
| `MaxRequestsPerConn` | 2 | Requêtes servies par connexion |
| `DisableKeepalive` | actif | Pas de connexion persistante : le flux infini joue ce rôle |
| `GetOnly` | actif | Seul `GET` est accepté |

Aucun `WriteTimeout` ni `IdleTimeout` n'est défini : un client qui lit très lentement reste piégé aussi longtemps qu'il le souhaite, ce qui est le but.

## La limite par IP derrière un reverse proxy

`MaxConnsPerIP` compte les connexions par **adresse TCP du pair**, pas par la valeur de l'en-tête `real_ip_header`. Derrière un reverse proxy (nginx, Nginx Proxy Manager...), tous les bots arrivent avec l'IP du proxy : la limite par défaut de 10 plafonne alors **toute l'instance** à 10 bots piégés en même temps, et les suivants reçoivent 429 au lieu d'être piégés. Chaque refus laisse une ligne de niveau debug `The number of connections from <ip> exceeds MaxConnsPerIP=<n>` : compter ces lignes dans les logs indique si la limite était atteinte. Voir [Journalisation](journalisation.md).

Jusqu'à ce fork la valeur était codée en dur à 10. Elle est maintenant configurable (`performance.max_conns_per_ip`, variable `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP`), avec la même valeur par défaut, y compris avec `-c` quand la clé est absente. Une valeur de quelques dizaines à centaines convient derrière un proxy ; `0` supprime toute limite. En mode socket Unix la limite ne s'applique pas du tout : voir [Serveur HTTP](../architecture/serveur-http-et-routage.md).

## Ce qui a été mesuré

Mesures indicatives, faites sur Docker Desktop (Windows) avec des clients `curl` bridés en débit, sans valeur de référence absolue :

- 30 connexions simultanées depuis une même IP avec la limite par défaut : 10 reçoivent 200, 20 reçoivent 429.
- 300 clients lents (2 Kio/s) avec la limite à 0 : tous acceptés. La mémoire du conteneur est d'environ 175 à 185 Mio, mais le processus HellPot lui-même n'occupe qu'environ 36 Mio : le reste est surtout constitué de tampons réseau du noyau, qui grossissent parce que les clients lisent lentement. La charge CPU reste négligeable (le débit est limité par la contre-pression TCP, pas par la génération de texte).

## Pourquoi le pool de buffers a disparu

L'ancien code gardait un `sync.Pool` de buffers de 100 Kio destinés à `io.CopyBuffer`. Il n'était jamais utilisé : la destination est un `*bufio.Writer`, qui implémente `io.ReaderFrom`, donc la copie ne passait pas par le buffer fourni. Le pool est supprimé ; le paramètre `buffsize` de `NewHeffalump` est conservé mais ignoré. La mesure ci-dessus a été faite avec et sans pool : aucune différence de mémoire (environ 36 Mio de RSS dans les deux cas). C'est un nettoyage, pas une optimisation.
