# Tmux Warp

Tmux Warp allows you to switch to a tmux session using the paths in your warp
config.

## The Why

My usual workflow to get start working on a project is:

```bash
tmux
wd warp-point
```

And then rename the session name to `warp-point`.

This project aims at reducing this to a single step:

```bash
twd start warp-point
```

## Usage

As already hinted at in the previous section, there is just a single command
provided currently, i.e., `start` (`s` can also be used).

```bash
twd start warp-point
```

The warp point must exist in your `.warprc` file; without which the command
returns an error.

The `start` command will work, both, inside and outside of a tmux session. If a
tmux server doesn't already exist, it creates that as well.

## Caveats

This project is designed to work for my specific usage pattern which is a
combination of the `wd` and `tmux` commands. This project not intended to
replace the `wd` command by any means. I do plan on adding a command to create a
new warp point and start a new session.

The errors reported by this project may be a little vague and will have some
rough edges.
