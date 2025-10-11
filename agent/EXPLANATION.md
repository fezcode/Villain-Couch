
# Say Hello to Your Friendly Neighborhood VLC Tracker Agent - Villain Couch!

Welcome to the official explanation of the **VLC Tracker Agent**! This little buddy is here to make your media-watching experience smoother and more enjoyable. Think of it as a personal assistant for your VLC Media Player.

## What does it do?

At its core, the agent does two main things:

1.  **Tracks your progress:** It remembers where you left off in your movies and TV shows, so you can easily resume watching later. No more scrubbing through a video to find your spot!
2.  **Helps you find what's next:** It can automatically find and play the next episode of a TV show for you, so you can binge-watch without interruption.

## Features and Flags

Here's a breakdown of all the cool things the agent can do, and how to use them:

### Adding a Workspace (`--ws`)

This is the first thing you'll want to do. A "workspace" is just a fancy word for a folder on your computer where you keep your media files. To add a workspace, you'll use the `--ws` flag, like this:

```bash
villain_couch --ws "/path/to/your/media/folder"
```

The agent will scan this folder and all its subfolders for media files. You can add as many workspaces as you want!

### Playing a File (`--file`)

If you want to play a specific file, you can use the `--file` flag:

```bash
villain_couch --file "/path/to/your/media/folder/movie.mp4"
```

The agent will start playing the file in VLC and will automatically start tracking your progress.

### Finding the Next Episode (`--find-next`)

This is where the magic happens! If you're watching a TV show, you can use the `--find-next` flag to tell the agent to automatically find and play the next episode when the current one finishes.

```bash
villain_couch --find-next
```

The agent is smart enough to figure out the next episode based on the filename. For example, if you're watching "My Awesome Show - S01E02.mkv", it will look for "My Awesome Show - S01E03.mkv" and play it for you.

### Verbose Mode (`--verbose`)

If you're a power user and want to see what the agent is doing under the hood, you can use the `--verbose` flag. This will make the agent print out more detailed logs.

```bash
villain_couch --verbose
```

## How it Works

The agent communicates with VLC Media Player in the background to get information about what you're watching. It then saves this information to a local database. When you start the agent again, it will automatically look up your progress and resume playback from where you left off.

We hope you enjoy using the Villain Couch! If you have any questions or feedback, please don't hesitate to reach out.
