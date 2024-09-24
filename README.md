# Diplomabot

Diplomabot is a Go-based bot that notifies me when my university diploma is ready for pickup. It automates the process of checking my diploma status on the IGCE website and sends me a Telegram message when the diploma is issued and ready to be
picked up.

## Overview

The bot can be described in five parts:

1. **Access the IGCE Website**
  The bot accesses the [IGCE Undergraduate Office's web page](https://igce.rc.unesp.br/index.php#!/instituicao/diretoria-tecnica-academica/graduacao/), which uses xajax for client-side rendering. It performs an http POST request that returns an XML containing the page's HTML.

2. **Unmarshal XML**
  After receiving the XML response, the bot unmarshals it to extract the HTML content from the body of the page.

3. **Scraping for Diploma Link**
  Using the HTML, the bot scrapes a link for a PDF document that lists the students names with issued diplomas.

4. **Download and Read PDF**
  The bot downloads the PDF and extract the text from the PDF using an external API.

5. **Notification via Telegram**
  If your name appears in the extracted text, the bot sends a notification message to you via Telegram, informing you that your diploma is ready for being picked up.

## Running Locally

### Clone the repository

  ```bash
  git clone https://github.com/vitorxfs/diplomabot
  cd diplomabot
  ```

### Set the following environment variables in a `.env` file

- `TELEGRAM_BOT_ID`: Your Telegram bot ID.
- `TELEGRAM_CHAT_ID`: Your Telegram chat ID for receiving notifications.
- `UTILS_TOKEN`: Token for the external API used to read PDF text.
- `UTILS_BASE_URL`: Base URL for the external API.

### Run as development mode

  ```sh
    go run .
  ```

### Compile and run

  ```sh
    go build -o ./.out/diplomabot
  ```

  ```sh
    ./.out/diplomabot
  ```

### Setting up the scheduler in ubuntu

1. Run `crontab -e` to open crontabs file;
2. Add your cron job in a new line on the file, adding your environment variables:

```txt
0 8 * * * TELEGRAM_BOT_ID=<bot_id> TELEGRAM_CHAT_ID=<chatId> UTILS_TOKEN=<utils_token> UTILS_BASE_URL=<utils_base_url> /path/to/binary/diplomabot >/path/to/logfile.txt 2>&1
```
