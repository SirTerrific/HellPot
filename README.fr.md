<div align="center">
  <img src="https://tcp.ac/i/00ctL.gif" alt="HellPot"/>

[![Vibe Check](https://github.com/SirTerrific/HellPot/actions/workflows/go.yml/badge.svg)](https://github.com/SirTerrific/HellPot/actions/workflows/go.yml) [![Docker (GHCR)](https://github.com/SirTerrific/HellPot/actions/workflows/docker.yml/badge.svg)](https://github.com/SirTerrific/HellPot/actions/workflows/docker.yml) [![GoDoc](https://godoc.org/github.com/SirTerrific/HellPot?status.svg)](https://godoc.org/github.com/SirTerrific/HellPot) [![Go Report Card](https://goreportcard.com/badge/github.com/SirTerrific/HellPot)](https://goreportcard.com/report/github.com/SirTerrific/HellPot) [![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](README.md) | **Français**

</div>

## À propos de ce fork

Ce dépôt est un fork maintenu de [yunginnanet/HellPot](https://github.com/yunginnanet/HellPot), dont le projet d'origine n'est plus maintenu.

L'objectif est un remplacement direct : **même comportement, même fichier de configuration, même format de logs.** Les changements se limitent à des dépendances à jour, un build rafraîchi, une image Docker publiée, quelques petites corrections et un nouveau paramètre optionnel. La liste exacte se trouve dans [CHANGELOG.md](CHANGELOG.md) (en anglais).

Si vous construisiez l'image d'origine localement, vous pouvez passer à l'image publiée sans modifier votre configuration ni votre chaîne de logs : voir [Méthode Docker](#méthode-docker).

## Résumé

HellPot est un honeypot sans fin, basé sur [Heffalump](https://github.com/carlmjohnson/heffalump), qui envoie les robots HTTP indisciplinés en enfer.

Il utilise notamment un [fichier de configuration toml](https://github.com/knadh/koanf), produit des [logs JSON](https://github.com/rs/zerolog) et offre des performances nettement supérieures.

## Conséquences funestes

Les clients (des robots, on l'espère) qui ignorent `robots.txt` et se connectent à votre instance de HellPot **subiront des conséquences éternelles**.

HellPot envoie un flux infini de données _juste assez proches_ d'un vrai site web pour que le client s'attarde jusqu'à ce que son âme soit déchirée et qu'il cesse d'exister.

Sous le capot de cette souffrance éternelle, un moteur de chaînes de Markov envoie au client des morceaux de [La Naissance de la tragédie (Hellénisme et Pessimisme)](https://www.gutenberg.org/files/51356/51356-h/51356-h.htm) de Friedrich Nietzsche, grâce à [fasthttp](https://github.com/valyala/fasthttp). Le texte source est en anglais.

## Installation

| Méthode | Commande / emplacement |
| --- | --- |
| Docker (multi-architecture `linux/amd64` + `linux/arm64`) | `docker pull ghcr.io/sirterrific/hellpot:latest` |
| Binaires compilés (Linux, Windows, macOS, FreeBSD) | [Releases GitHub](https://github.com/SirTerrific/HellPot/releases/latest) |
| Depuis les sources | voir [Compiler depuis les sources](#compiler-depuis-les-sources) |

## Compiler depuis les sources

HellPot nécessite **Go 1.26 ou plus récent** (voir la ligne `go` de [go.mod](go.mod)). Il utilise les [modules Go](https://go.dev/blog/using-go-modules), ce qui rend la compilation très simple avec une installation Go standard. Un Makefile GNU est fourni.

1 ) `git clone https://github.com/SirTerrific/HellPot`

2 ) `cd HellPot`

3 ) `make`

4 ) _Mesurez les conséquences funestes potentielles de vos actes._

Les cibles du Makefile sont `deps` (`go mod tidy`), `check` (`go vet`), `build`, `run` et `format` ; `make` seul exécute `deps check build`. `make build` écrit le binaire `HellPot` dans le répertoire courant.

Pour installer directement depuis le chemin du module :

```
go install github.com/SirTerrific/HellPot/cmd/HellPot@latest
```

## Utilisation

### Méthode YOLO :

Si aucun fichier de configuration n'existe, HellPot tente de placer sa configuration par défaut dans **$HOME/.config/HellPot/config.toml** (le répertoire de configuration utilisateur de votre système). Cela permet aux âmes imprudentes de faire pleuvoir le feu **_immédiatement_** :

1 ) Téléchargez une [version compilée](https://github.com/SirTerrific/HellPot/releases/latest)

2 ) Lancez le binaire et envoyez aussitôt les clients directement en enfer.

---

### Méthode raisonnable :

1 ) Configurez le serveur web comme reverse proxy (voir plus bas)

2 ) `./HellPot --genconfig`

3 ) Modifiez le `config.toml` généré à votre convenance.

4 ) Réfléchissez à la capacité de votre serveur à encaisser les valeurs de performance choisies.

5 ) `./HellPot -c config.toml`

---

### Méthode Docker

Une image multi-architecture est publiée sur le GitHub Container Registry. Elle est reconstruite à chaque push sur `main`, à chaque tag `v*` et chaque lundi (pour récupérer les correctifs de sécurité de Go et de l'image de base).

```
docker pull ghcr.io/sirterrific/hellpot:latest
```

| Tag | Publié quand |
| --- | --- |
| `latest` | push sur `main`, et reconstruction hebdomadaire |
| `sha-<commit court>` | chaque push sur `main` ou sur un tag `v*`, et reconstruction hebdomadaire |
| `<version>` et `<majeure>.<mineure>` | un tag `v*` est poussé (par exemple `v1.2.3` donne `1.2.3` et `1.2`) |

La structure de l'image est identique à celle du `Dockerfile` d'origine :

| Élément | Valeur |
| --- | --- |
| Binaire | `/app` |
| Point d'entrée | `/app -c /config` |
| Configuration par défaut | `/config` (copie de [docker_config.toml](docker_config.toml)) |
| Répertoire des logs | `/logs/` |
| Port | `8080` |
| Image de base | `gcr.io/distroless/static-debian13` (pas de shell) |
| Utilisateur | root (uid 0) |

La configuration par défaut écoute sur `0.0.0.0:8080`, active `catchall` (tous les chemins GET reçoivent le piège, aucun `robots.txt` n'est servi), fixe l'en-tête `Server` à `nginx`, met sur liste noire les user agents contenant `Cloudflare-Traffic-Manager` ou `curl`, lit l'IP du client dans `X-Real-IP`, et écrit ses logs JSON dans `/logs/` avec une sortie console sans couleur.

Démarrage rapide :

```
docker run -d --name hellpot -p 8080:8080 -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest
```

Utilisez votre propre configuration en montant un fichier sur `/config` :

```
docker run -d --name hellpot -p 8080:8080 -v ./config.toml:/config:ro -v ./logs:/logs ghcr.io/sirterrific/hellpot:latest
```

Vous pouvez aussi surcharger des clés une à une avec les [variables d'environnement](#variables-denvironnement). Exemple `docker compose`, derrière un reverse proxy qui partage le réseau externe `proxy` :

```yaml
services:
  hellpot:
    image: ghcr.io/sirterrific/hellpot:latest
    container_name: hellpot
    restart: unless-stopped
    expose:
      - 8080
    volumes:
      - ./logs:/logs
    environment:
      - TZ=America/Toronto
      # tous les clients partagent l'IP du reverse proxy : relever la limite de connexions par IP
      - HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP=100
    # durcissement optionnel, vérifié avec cette image
    read_only: true
    cap_drop: [ALL]
    security_opt:
      - no-new-privileges:true
    networks:
      - proxy

networks:
  proxy:
    external: true
```

`TZ` définit le fuseau horaire utilisé dans le nom du fichier de log et dans le champ `time` de chaque ligne de log.

> [!NOTE]
> Si vous passez vous-même `-c` (`docker run ... <image> -c autre.toml`), les arguments sont ajoutés au point d'entrée : `/config` est donc chargé en premier, puis votre fichier est fusionné par-dessus.

---

666 ) 𝙏͘͝𝙝̓̓͛𝙚͑̈́̀ 𝙨͆͠͝𝙠͑̾͌𝙮̽͌͆ 𝙞̓̔̔𝙨͒͐͝ 𝙛͑̈́̚𝙖͛͒𝙡͑͆̽𝙡̾̚̚𝙞͋̒̒𝙣̾͛͝𝙜͒̒̀.́̔͝​

## Ligne de commande

| Option | Effet |
| --- | --- |
| `-c`, `--config <fichier>` | Utilise ce fichier de configuration. **Les valeurs par défaut intégrées ne sont alors plus appliquées**, voir la note plus bas. |
| `-v`, `--debug` | Force le niveau debug |
| `-vv`, `--trace` | Force le niveau trace |
| `--nocolor` | Désactive les couleurs et la bannière |
| `--banner` | Affiche la bannière et la version, puis quitte |
| `--genconfig` | Écrit la configuration par défaut dans `./config.toml`, puis quitte |
| `-h` | Affiche l'aide (seulement si la sortie standard est un terminal), puis quitte |

> [!NOTE]
> Le texte d'aide mentionne `--help`, mais seul `-h` est reconnu : `--help` est ignoré et le serveur démarre. Le texte d'aide indique aussi que `--genconfig` écrit `HellPot.toml`, alors que le fichier écrit est `config.toml`.

La version affichée par `--banner` est le tag git à partir duquel le binaire a été compilé (par exemple `1.2.3` pour le tag `v1.2.3`). Un build sans tag affiche le hash court du commit, ou `dev`.

## Référence de configuration

Sans option `-c`, HellPot charge `/etc/HellPot/config.toml` s'il existe (sauf sous Windows), sinon `config.toml` dans votre répertoire de configuration utilisateur, créé à partir des valeurs par défaut s'il manque. Les valeurs par défaut intégrées complètent toutes les clés absentes de votre fichier.

> [!IMPORTANT]
> Avec `-c`, **les valeurs par défaut intégrées ne sont pas appliquées.** Toute clé omise prend sa valeur zéro (chaîne vide, `false`, `0`, liste vide). Par exemple, un fichier ne contenant que les paramètres d'écoute `[http]` et aucune section `[http.router]` n'enregistre aucune route : tous les chemins répondent `404` et il n'y a pas de `robots.txt`. Définissez explicitement `catchall = true` ou `paths`. La seule exception est `performance.max_conns_per_ip`, qui reste à `10` quand elle est absente.

### Variables d'environnement

Les valeurs de configuration peuvent être surchargées par des variables d'environnement préfixées par `HELLPOT_`. Utilisez un seul tiret bas pour chaque niveau `.` et **deux tirets bas** pour chaque tiret bas à l'intérieur d'un nom de clé :

| Clé de configuration | Variable d'environnement |
| --- | --- |
| `http.bind_addr` | `HELLPOT_HTTP_BIND__ADDR` |
| `http.router.catchall` | `HELLPOT_HTTP_ROUTER_CATCHALL` |
| `logger.docker_logging` | `HELLPOT_LOGGER_DOCKER__LOGGING` |
| `performance.max_conns_per_ip` | `HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP` |

### Référence

Le bloc ci-dessous correspond à ce que produit `./HellPot --genconfig` (commentaires ajoutés). Quand la valeur par défaut de l'image Docker diffère, c'est indiqué.

```toml
[deception]
  # Utilisé comme en-tête HTTP "Server". Un reverse proxy peut le masquer.
  server_name = "nginx"

[http]
  # Écoute TCP (par défaut)
  bind_addr = "127.0.0.1"   # l'image Docker utilise "0.0.0.0"
  bind_port = "8080"

  # nom de l'en-tête contenant la vraie IP du client, pour les déploiements derrière reverse proxy
  real_ip_header = 'X-Real-IP'

  # liste de chaînes de user agent en liste noire (sensible à la casse).
  # les clients dont le user agent contient l'une de ces chaînes reçoivent "Not found" pour toute requête.
  # défaut : ["Cloudflare-Traffic-Manager"] ; l'image Docker ajoute "curl"
  uagent_string_blacklist = ["Cloudflare-Traffic-Manager", "curl"]

  # Écoute sur socket Unix (remplace l'écoute TCP). Non pris en charge sous Windows.
  unix_socket_path = "/var/run/hellpot"
  unix_socket_permissions = "0666"
  use_unix_socket = false

  [http.router]
    # À true, toutes les requêtes GET correspondent. Force makerobots = false.
    catchall = false          # l'image Docker utilise true
    # À false, le gestionnaire de robots.txt n'est pas créé.
    makerobots = true
    # Des gestionnaires sont créés pour ces chemins, ainsi que des entrées dans robots.txt. Valide seulement si catchall = false.
    paths = ["wp-login.php", "wp-login"]

[logger]
  # verbeux (-v)
  debug = true
  # très verbeux (-vv)
  trace = false
  # les fichiers de log JSON sont écrits dans ce répertoire.
  # défaut si vide : $HOME/.local/share/HellPot/logs ; l'image Docker utilise "/logs/"
  directory = "/home/kayos/.local/share/HellPot/logs/"
  # désactive toute couleur dans la console. Sous Windows, vaut true par défaut.
  nocolor = false
  # utilise la date courante dans le nom des nouveaux fichiers de log.
  use_date_filename = true
  # écrit le JSON uniquement sur stdout (pas de fichier de log), implique nocolor. Prévu pour `docker logs` / les collecteurs de logs.
  docker_logging = false
  # format d'heure Go utilisé pour l'horodatage de la console (pas dans les logs JSON)
  console_time_format = "3:04PM"

[performance]
  # max_workers n'est valide que si restrict_concurrency est true
  max_workers = 256
  restrict_concurrency = false
  # nombre maximal de connexions simultanées depuis une même IP ; au-delà, le client reçoit HTTP 429.
  # Derrière un reverse proxy, tous les clients ont l'IP du proxy : cette valeur plafonne donc toute l'instance.
  # 0 = illimité. Ajouté dans ce fork ; la valeur par défaut, 10, est celle qui était codée en dur.
  max_conns_per_ip = 10
```

`catchall = true` prend le pas sur `makerobots` : aucun gestionnaire de `robots.txt` n'est créé. Avec `catchall = false`, une requête vers un chemin absent de `paths` (et différent de `/robots.txt`) reçoit `404`.

## Comportement

- Seul `GET` est servi. Les autres méthodes sont rejetées par le serveur HTTP.
- Une requête qui correspond reçoit `<html><body>` suivi d'un flux sans fin généré par une chaîne de Markov. La réponse ne se termine que lorsque le client part.
- Le délai de lecture de la requête est de 5 s, le corps de requête maximal est de 1 Mio, le keep-alive est désactivé et une connexion sert au plus 2 requêtes.
- Les user agents contenant une chaîne de la liste noire reçoivent `404 Not found` et ne sont pas piégés.
- `robots.txt` (s'il est activé) est généré à partir de `paths`, une ligne `Disallow:` par chemin : les robots qui l'ignorent sont ceux qui se font piéger.
- L'en-tête `Server` vaut `deception.server_name`.
- Le **mode socket Unix** (`use_unix_socket = true`) est servi par un serveur fasthttp par défaut : les délais, `max_conns_per_ip`, le mode GET seulement et l'en-tête `Server` issu de `deception.server_name` décrits ici ne s'y appliquent **pas** (l'en-tête `Server` vaut alors `fasthttp`). Ce comportement vient du projet d'origine et n'a pas été modifié.
- L'adresse de client journalisée est la valeur de `real_ip_header` quand elle est présente, sinon l'adresse TCP du pair.

## Logs

HellPot écrit un objet JSON par ligne (zerolog) et, sauf si `docker_logging` est activé, affiche aussi une version lisible sur la console. Le fichier de log est `<directory>/HellPot[_<date et heure de démarrage>].log`, par exemple `HellPot_03_Oct_26_20-56_EDT.log`, nommé une seule fois au démarrage du processus. Avec `docker_logging = true`, le JSON va sur stdout et aucun fichier n'est écrit.

Lignes par requête :

```json
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","time":"2026-10-03T22:00:36Z","message":"NEW"}
{"level":"info","USERAGENT":"Mozilla/5.0","REMOTE_ADDR":"1.2.3.4","URL":"/wp-login.php","BYTES":179549,"DURATION":15.914698,"time":"2026-10-03T22:00:36Z","message":"FINISH"}
```

| `message` | Niveau | Quand | Champs supplémentaires |
| --- | --- | --- | --- |
| `NEW` | info | un client a été envoyé en enfer | `USERAGENT`, `REMOTE_ADDR`, `URL` |
| `FINISH` | info | le client est parti | `BYTES` envoyés, `DURATION` en millisecondes |
| `END_ON_ERR` | trace | l'écriture a échoué (en général le client s'est déconnecté) | `error` |
| `Ignoring useragent` | trace | un user agent en liste noire a été refusé | |
| `SERVE_ROBOTS` | debug | `robots.txt` a été servi | `PATHS` |
| `The number of connections from <ip> exceeds MaxConnsPerIP=<n>` | debug | un client a été refusé avec HTTP 429 | |

`caller` est ajouté aux lignes de requête quand `trace` est activé. Les lignes de démarrage (`config`, `logger`, `Listening and serving HTTP...`) indiquent le fichier de configuration, le fichier de log et l'adresse d'écoute. Les messages de log restent en anglais, comme dans le projet d'origine.

> [!NOTE]
> Particularité connue, conservée volontairement pour que le contenu des logs reste identique à celui du projet d'origine : les lignes debug et trace sont actuellement écrites même avec `debug = false` et `trace = false`. Ces deux paramètres ne peuvent que relever le niveau.

## Exemples de configuration de reverse proxy

#### nginx

<details>
  <summary>nginx</summary>

```nginx
location '/robots.txt' {
	proxy_set_header Host $host;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_pass http://127.0.0.1:8080$request_uri;
}

location '/wp-login.php' {
	proxy_set_header Host $host;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_pass http://127.0.0.1:8080$request_uri;
}
```

</details>

#### Apache

<details>
  <summary>apache (mod_proxy + mod_proxy_http)</summary>

Toutes les URL inexistantes sont envoyées par reverse proxy vers une instance de HellPot sur localhost, configurée en catchall. Le trafic servi par HellPot est limité à 5 Kio/s.

- Créez votre robots.txt et votre contenu habituel. Créez aussi le répertoire et les fichiers factices des Errordocument (les fichiers peuvent être vides). Dans l'exemple, le répertoire est "/content/"
- Une requête vers une URL qui a un gestionnaire existant (par ex. un fichier) est traitée par apache
- Les requêtes vers des URL inexistantes provoquent une erreur HTTP 404, dont le contenu est servi par HellPot
- Les URL sous "/.well-known/" sont exclues.

```apache
<VirtualHost yourserver>
    ErrorDocument 400 "/content/400"
    ErrorDocument 403 "/content/403"
    ErrorDocument 404 "/content/404"
    ErrorDocument 500 "/content/405"
    <Directory "$wwwroot/.well-known/">
        ErrorDocument 400 default
        ErrorDocument 403 default
        ErrorDocument 404 default
        ErrorDocument 500 default
    </Directory>
    /* HTTP Honeypot / HellPot (need mod_proxy, mod_proxy_http) */
    ProxyPreserveHost	on
    ProxyPass         "/content/" "http://localhost:8080/"
    ProxyPassReverse  "/content/" "http://localhost:8080/"

    /* Rate Limit config, need mod_ratelimit */
    <Location "/content/">
        SetOutputFilter RATE_LIMIT
        SetEnv rate-limit 5
    </Location>

    /* Remaining config */

</VirtualHost>
```

</details>

## Notes de sécurité

- **`real_ip_header` est cru tel quel.** HellPot ne vérifie pas qui l'a envoyé. N'exposez HellPot que derrière un reverse proxy qui écrase cet en-tête (comme dans l'exemple nginx), ou gardez-le sur un réseau interne, sinon un client peut falsifier le `REMOTE_ADDR` qui se retrouve dans vos logs.
- L'image Docker s'exécute en root. Elle n'a besoin d'aucune capability ajoutée ni de système de fichiers inscriptible en dehors de `/logs` : les paramètres `read_only`, `cap_drop` et `no-new-privileges` de l'exemple compose ci-dessus sont vérifiés comme fonctionnels.
- La CI exécute `go vet`, `gosec`, les tests avec le détecteur de concurrence et `govulncheck` à chaque push (voir plus bas), et l'image Docker est reconstruite chaque semaine.

## Développement

- **CI** ([go.yml](.github/workflows/go.yml), « Vibe Check ») : à chaque push et sur les pull requests vers `main` : `go vet`, `gosec`, `go test -race`, `go build`, `govulncheck`. La version de Go vient de `go.mod`.
- **Docker** ([docker.yml](.github/workflows/docker.yml)) : construit `linux/amd64` et `linux/arm64` et publie sur `ghcr.io/sirterrific/hellpot` (les pull requests construisent seulement, sans publier). Le Dockerfile exécute `go vet` et `go test` avant la compilation, et accepte un argument de build `VERSION`.
- **Releases** ([release-command.yml](.github/workflows/release-command.yml)) : la création d'une release GitHub compile des binaires pour linux, windows, darwin et freebsd (386, amd64, arm64, sauf darwin/386 et windows/arm64), avec les sommes SHA-256.
- **Dependabot** vérifie les modules Go et les GitHub Actions chaque jour, et les images de base Docker chaque semaine.

Pour publier une version, poussez un tag comme `v1.2.3` : le workflow Docker publie les tags d'image versionnés. Créez ensuite la release GitHub à partir de ce tag pour obtenir les binaires.

Une base de connaissances générée sur le code est disponible dans le répertoire [openwiki](openwiki).

## Autres souffrances

- https://github.com/ginger51011/pandoras_pot
  - Un honeypot HTTP inspiré de HellPot, pour punir et éduquer les robots d'indexation indisciplinés, écrit en Rust (🚀)

## Crédits et licence

HellPot a été créé par [yung innanet](https://github.com/yunginnanet) et repose sur [Heffalump](https://github.com/carlmjohnson/heffalump) de Carl Johnson. Ce fork conserve la [licence MIT](LICENSE) d'origine.
