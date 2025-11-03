# LyriFlow

🎼 Get synchronized lyrics for your music

指定した再生位置の歌詞を表示するCLIツール

## Inspired

- [waybar-lyric](https://github.com/Nadim147c/waybar-lyric)

## インストール

```bash
go install github.com/arrow2nd/lyriflow@latest
```

または、リポジトリをクローンしてビルド：

```bash
git clone https://github.com/arrow2nd/lyriflow.git
cd lyriflow
go build
```

## 使い方

### コマンド一覧

#### `get` - 歌詞を取得

指定した再生位置の歌詞を取得して表示します

```bash
lyriflow get -t "曲名" -a "アーティスト" -A "アルバム" -p 30.5
```

**必須オプション：**

- `-t, --title`: 曲のタイトル
- `-a, --artist`: アーティスト名
- `-A, --album`: アルバム名
- `-p, --position`: 現在の再生位置（秒単位）

**オプション：**

- `--waybar`: waybar用のJSON形式で出力

#### `cache-purge` - キャッシュをクリア

保存された歌詞キャッシュを削除します

```bash
lyriflow cache-purge
```

#### `version` / `v` - バージョン表示

LyriFlowのバージョンを表示します

```bash
lyriflow version
```

## waybar連携

### JSON出力フォーマット

`--waybar`フラグを使用すると、以下の形式でJSON出力されます：

```json
{
  "text": "表示テキスト",
  "alt": "状態識別子",
  "tooltip": "ツールチップ（曲情報）",
  "class": "CSSクラス名"
}
```

### 状態別の出力

| 状態                 | `text`                | `alt`          | `class`        |
| -------------------- | --------------------- | -------------- | -------------- |
| 歌詞表示中           | 歌詞テキスト          | `playing`      | `lyrics`       |
| 間奏中               | `(instrumental)`      | `instrumental` | `instrumental` |
| 同期できる歌詞がない | `No lyrics available` | `no-lyrics`    | `no-lyrics`    |
| 歌詞が見つからない   | `Lyrics not found`    | `not-found`    | `not-found`    |

### waybar設定例

#### config

```json
{
  "custom/lyrics": {
    "return-type": "json",
    "format": "{icon} {0}",
    "hide-empty-text": true,
    "exec": "lyriflow get -t \"$(playerctl metadata title)\" -a \"$(playerctl metadata artist)\" -A \"$(playerctl metadata album)\" -p $(playerctl position) --waybar",
    "interval": 1,
    "on-click": "playerctl play-pause"
  }
}
```

#### style.css

```css
#custom-lyrics.lyrics {
  color: #a6e3a1;
}

#custom-lyrics.instrumental {
  color: #89b4fa;
}

#custom-lyrics.no-lyrics {
  color: #f38ba8;
}

#custom-lyrics.not-found {
  color: #fab387;
}
```
