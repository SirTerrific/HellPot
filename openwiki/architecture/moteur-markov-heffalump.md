---
type: architecture
title: Moteur de chaînes de Markov (heffalump)
description: Comment le package heffalump fabrique le flux de texte infini envoyé aux robots, de la décompression du texte source à l'écriture dans la réponse HTTP, et comment il s'arrête quand le client part.
tags: [heffalump, markov, streaming, tarpit]
verified:
  - by: openwiki/0.5.1
    at: 2026-10-04T19:20:53.104Z
sources:
  - id: openwiki-source-a97afb7f0fe6ba7d781c5467
    resource: repo://heffalump/heffalump_test.go
  - id: openwiki-source-6225717d77c0df8b3a43d4df
    resource: repo://heffalump/heffalump.go
  - id: openwiki-source-db1e1e64f92e7e42faa4177e
    resource: repo://heffalump/markov.go
  - id: openwiki-source-e689a4a46f2ebef989178800
    resource: repo://internal/http/router.go
generated: { by: "claude-code", at: "2026-10-04T02:02:24.687Z" }
---

# Moteur de chaînes de Markov (heffalump)

Le package `heffalump` produit le contenu du piège : un flux de texte pseudo-aléatoire, sans fin, qui ressemble à du HTML. Il est issu du projet Heffalump de Carl Johnson (voir l'en-tête de [heffalump.go](../../heffalump/heffalump.go)). Il ne connaît ni HTTP ni la configuration, à l'exception du logger partagé. Son seul client est le gestionnaire de requêtes décrit dans [Serveur HTTP, routage et piège](serveur-http-et-routage.md).

## Initialisation : table construite une fois au démarrage

Au chargement du package, `init()` décompresse le texte source embarqué (`srcGz`, dans `src.go`) avec `squish.UnpackStr`, puis construit `DefaultMarkovMap` et `DefaultHeffalump`. Le processus **panique** si la décompression échoue ou si le texte obtenu est vide : il n'y a pas de mode dégradé. La table est donc entièrement en mémoire avant que le serveur n'écoute, et elle est ensuite en lecture seule.

Le texte source est *La Naissance de la tragédie* de Nietzsche (en anglais), comme l'indique le [README](../../README.md).

## Structure de données

- `MarkovMap` associe une **paire** de tokens consécutifs (`tokenPair`) à la liste de tous les tokens qui ont suivi cette paire dans le texte. Les doublons sont conservés : un suffixe fréquent apparaît plusieurs fois, ce qui le rend plus probable.
- `Fill` lit le texte avec un `bufio.Scanner` dont la fonction de découpage est `ScanHTML`. Celle-ci renvoie soit un mot (délimité par des espaces ou par un `<`), soit une balise HTML entière de `<` à `>`. Les balises sont donc traitées comme des tokens, ce qui permet au texte généré de contenir de vraies balises.
- La paire de départ est `("", "")`.

## Génération

- `Get(w1, w2)` choisit un suffixe au hasard parmi ceux de la paire, avec `math/rand`. Ce n'est pas un générateur cryptographique, ce qui est assumé (commentaire `#nosec` pour gosec G404). Si la paire est inconnue, il renvoie la chaîne vide.
- `Read(p)` implémente `io.Reader` : il enchaîne les appels à `Get`, copie chaque mot suivi d'un saut de ligne dans `p`, et s'arrête quand le mot suivant ne tiendrait plus. Il ne renvoie **jamais** d'erreur ni `io.EOF` : la source est inépuisable.

Quand la chaîne atteint la fin du texte, `Get` renvoie des chaînes vides. D'après le code, l'état redevient alors `("", "")`, c'est-à-dire la paire de départ, et la génération reprend depuis le début du texte.

## Écriture dans la réponse : `WriteHell`

`WriteHell(bw *bufio.Writer)` écrit d'abord l'en-tête `<html>\n<body>\n`, puis appelle `io.Copy(bw, h.mm)`. Comme `MarkovMap.Read` ne finit jamais, cet appel ne rend la main que lorsqu'une **écriture échoue**, c'est-à-dire en pratique quand le client s'est déconnecté. Il renvoie alors le nombre d'octets écrits et une erreur `nil` : l'erreur de copie est volontairement ignorée.

La fin du flux est détectée un cran plus loin. La boucle du routeur rappelle `WriteHell` tant qu'il ne renvoie pas d'erreur ; à l'appel suivant, l'écriture de l'en-tête échoue (le `bufio.Writer` conserve son erreur), `WriteHell` renvoie cette erreur, et la boucle s'arrête en journalisant `END_ON_ERR` puis `FINISH`. Voir [Journalisation](../operations/journalisation.md).

Une panique survenue pendant l'écriture est interceptée par un `recover` : elle est journalisée au niveau `error` avec le message `panic recovered!`, et ne fait pas tomber le processus.

## Pas de pool de buffers

Une version antérieure gardait un `sync.Pool` de buffers de 100 Kio passé à `io.CopyBuffer`. Il n'était jamais utilisé : `*bufio.Writer` implémente `io.ReaderFrom`, donc la copie passait par son propre mécanisme. Le pool a été supprimé et le paramètre `buffsize` de `NewHeffalump` est conservé pour compatibilité de signature mais ignoré. Aucune différence de sortie ni de logs. Voir [Limites de connexions et performance](../operations/limites-et-performance.md).

## Vérification

`heffalump_test.go` contient l'unique test du projet : il écrit dans un `Writer` qui échoue après 64 Kio (simulation d'un client qui raccroche), rejoue la boucle du routeur, et vérifie que (1) la boucle se termine, (2) des octets ont été écrits, (3) la sortie commence par l'en-tête HTML. Il est exécuté par le Dockerfile et par la CI ([CI, sécurité et releases](../operations/ci-securite-et-releases.md)).
