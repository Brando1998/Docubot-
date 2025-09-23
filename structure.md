.
├── api
│   ├── cmd
│   │   └── api
│   ├── config
│   │   └── config.go
│   ├── controllers
│   │   ├── auth.go
│   │   ├── automation.go
│   │   ├── conversation.go
│   │   ├── document.go
│   │   ├── health.go
│   │   ├── user.go
│   │   ├── user_test.go
│   │   ├── websocket_hub.go
│   │   └── whatsapp.go
│   ├── databases
│   │   ├── mongo.go
│   │   └── postgresql.go
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── .env
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   ├── middleware
│   │   ├── middleware.go
│   │   └── paseto.go
│   ├── mocks
│   │   └── client_repository_mock.go
│   ├── models
│   │   ├── bot.go
│   │   ├── client.go
│   │   ├── messages.go
│   │   ├── paseto.go
│   │   ├── system_user.go
│   │   └── whatsapp.go
│   ├── repositories
│   │   ├── bot_repository.go
│   │   ├── client_repository.go
│   │   └── conversation_repository.go
│   ├── routes
│   │   └── routes.go
│   └── services
├── baileys-ws
│   ├── auth
│   │   ├── app-state-sync-key-AAAAAO__I.json
│   │   ├── app-state-sync-version-critical_block.json
│   │   ├── app-state-sync-version-critical_unblock_low.json
│   │   ├── app-state-sync-version-regular_high.json
│   │   ├── app-state-sync-version-regular.json
│   │   ├── app-state-sync-version-regular_low.json
│   │   └── creds.json
│   ├── .env.example
│   ├── package.json
│   ├── package-lock.json
│   ├── src
│   │   ├── handlers
│   │   ├── index.ts
│   │   ├── sessions
│   │   └── websocket
│   └── tsconfig.json
├── docker
│   └── api.Dockerfile
├── docker-compose.yml
├── k8s
│   ├── configmaps
│   │   └── api-config.yaml
│   ├── deployments
│   │   ├── api-deployment.yaml
│   │   ├── mongo-deployment.yaml
│   │   └── postgres-deployment.yaml
│   ├── secrets
│   │   └── api-secret.yaml
│   └── services
│       ├── api-service.yaml
│       ├── mongo-service.yaml
│       └── postgres-service.yaml
├── Makefile
├── playwright-bot
│   └── .env.example
├── rasa-bot
│   ├── actions
│   │   ├── actions.py
│   │   ├── __init__.py
│   │   └── __pycache__
│   ├── config.yml
│   ├── credentials.yml
│   ├── data
│   │   ├── nlu.yml
│   │   ├── rules.yml
│   │   └── stories.yml
│   ├── domain.yml
│   ├── endpoints.yml
│   ├── .env.example
│   ├── __init__.py
│   ├── models
│   │   ├── 20250808-102648-proper-fluid.tar.gz
│   │   ├── 20250815-213007-bone-milk.tar.gz
│   │   └── 20250815-214034-dichotomic-faucet.tar.gz
│   ├── .rasa
│   │   └── cache
│   ├── readme.md
│   └── requirements.txt
├── README.md
└── structure.md

34 directories, 74 files
