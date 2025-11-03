#!/bin/bash

# 再生中のメタデータを取得
title=$(playerctl metadata title 2>/dev/null)
artist=$(playerctl metadata artist 2>/dev/null | cut -d',' -f1 | xargs)

if [ -z "$title" ] || [ -z "$artist" ]; then
    echo '{"text":"No track playing","alt":"not-found","tooltip":"No track is currently playing","class":"not-found"}'
    exit 0
fi

album=$(playerctl metadata album 2>/dev/null)

# albumが取得できない場合は空文字列を使用
if [ -z "$album" ]; then
    album=""
fi

position=$(playerctl position 2>/dev/null | awk '{print $1 + 0.5}')

lyriflow get --title "$title" --artist "$artist" --album "$album" --position "$position" --waybar
