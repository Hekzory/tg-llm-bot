# Telegram LLM Bot

A Telegram bot that leverages Large Language Models (LLMs) through Ollama to provide AI-powered conversations. The system is built with a microservices architecture using Go, with separate services for handling Telegram interactions and model inference.

## Architecture

The project consists of two main services:

- **telegram-service**: Handles Telegram bot interactions and user management
- **model-service**: Manages LLM interactions through Ollama
- **shared**: Common utilities and models used by both services

### Key Features

- Microservices architecture for scalability and separation of concerns
- PostgreSQL database for persistent storage
- Docker containerization for easy deployment
- Support for multiple LLM models through Ollama
- Configurable through TOML files

## Prerequisites

- Docker and Docker Compose
- Go 1.24.1 or later (for development)
- NVIDIA GPU + drivers (optional, for GPU acceleration)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/Hekzory/tg-llm-bot.git
cd tg-llm-bot
```

2. Fill the configuration files with your own values and Telegram bot key:

- `config/model-service.toml`
- `config/telegram-service.toml`

3. Build and start the services:

For CPU-only usage:

```bash
docker compose up --build
```

or 

```bash
make up
```

For Nvidia GPU speed-up:

```bash
docker compose -f docker-compose-nvidia.yml up --build
```

4. Access the bot on Telegram:

- Search for `@your_bot_username` in Telegram


## Acknowledgments

- [Ollama](https://ollama.ai/) for providing the LLM serving infrastructure