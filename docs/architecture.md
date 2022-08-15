# Architecture

This document describes the high-level architecture of ggrok. If you want to familiarize yourself with the code base, you are just in the right place!

## Bird's Eye View

```mermaid

flowchart LR
    PublicNet <--http--> ggrok-server
    ggrok-server <-- websocket --> ggrok-client
    PublicNet x--x LocalServer
    LocalServer <--http--> ggrok-client

```

On the highest level, ggrok is a thing that exposes your local server to the internet, so that others can visit your local server easily.