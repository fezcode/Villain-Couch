# Villain Couch - VLC Tracker Agent

The Villain Couch VLC Tracker Agent is a lightweight, command-line application that seamlessly integrates with VLC Media Player to track your media playback progress. It runs in the background, monitoring VLC's status and saving the playback position of your media files. This allows you to resume watching your favorite movies and TV shows right where you left off, even after closing VLC or restarting your computer.

## Features

- **Automatic Playback Tracking:** Monitors VLC and automatically saves the playback progress of your media files.
- **Seamless Resumption:**  Lets you pick up right where you left off.
- **Graceful Shutdown:**  Ensures that your playback progress is saved even when you close the agent with `Ctrl+C`.
- **Easy Configuration:** Uses a simple `config.json` file for configuration.
- **Command-Line Interface:** Provides a simple and easy-to-use command-line interface.

## Installation

To build the agent, you need to have Go installed on your system.

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/villian-couch.git
   ```

2. **Navigate to the agent directory:**
   ```bash
   cd villian-couch/agent
   ```

3. **Build the agent:**
   ```bash
   go build -o villian-couch ./src
   ```

This will create a `villian-couch` executable in the `agent` directory.

## Configuration

The agent uses a `config.json` file for configuration. The first time you run the agent, it will automatically create a default `config.json` file in your user's configuration directory.

The default configuration looks like this:

```json
{
  "web_url": "http://localhost",
  "status_endpoint": "requests/status.json",
  "playlist_endpoint": "requests/playlist.json",
  "http_port": "9713",
  "extra_intf": "http",
  "http_password": "my_secret_password",
  "database_file_name": "storage.sqlite"
}

```

- `web_url`: The URL of the VLC web interface.
- `status_endpoint`: The endpoint for getting the VLC status.
- `playlist_endpoint`: The endpoint for getting the VLC playlist.
- `http_port`: The port for the VLC web interface.
- `extra_intf`: The extra interface to use for VLC.
- `http_password`: The password for the VLC web interface.
- `database_file_name`: The name of the database file where the playback progress will be stored.

## Usage

To run the agent, use the following command:

```bash
./villian-couch [flags]
```

### Flags

- `--verbose`: Enable verbose logging.
- `--file <media-file>`: Specify a media file to play.
- `--ws <directory>`: Adds given directory as workspace to find possible next episodes using `find-next` flag.
- `--find-next`: Try to find next episode in workspace.

### Examples

- **Start tracking and play a media file:**
  ```bash
  ./villian-couch --file /path/to/your/media.mp4
  ```

- **Enable verbose logging:**
  ```bash
  ./villian-couch --verbose --file /path/to/your/media.mp4
  ```

## Development

To contribute to the development of the agent, you can follow these steps:

1. **Fork the repository.**
2. **Create a new branch for your feature or bug fix.**
3. **Make your changes and commit them.**
4. **Push your changes to your fork.**
5. **Create a pull request.**

### Building for Development

During development, you can use the following command to build and run the agent:

```bash
go run ./src [flags]
```