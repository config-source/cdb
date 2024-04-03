#!/usr/bin/env sh

SUBCOMMAND="$1"

if [ -z "$SUBCOMMAND" ]; then
    SUBCOMMAND="server"
else
    shift
fi

if [ -n "$RUN_MIGRATIONS" ]; then
    /app/cdbd migrate
fi

/app/cdbd "$SUBCOMMAND" "$@"
