# `snipr`

[![Code Quality](https://img.shields.io/badge/Code_Quality-A+-brightgreen?style=for-the-badge&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/lucasepe/snipr)

**`snipr` — Execute Markdown code blocks with precision**

> Turn your Markdown documents into runnable scripts.

Run shell commands and Go code **directly from Markdown**, respecting dependencies, exporting outputs, and saving to files.


## 🚀 Features

- execute **runnable code blocks** in Markdown
- **topological sorting** of blocks based on `depends`
- export block output to **environment variables** (`export`) or **files** (`file`)
- skip blocks selectively with `skip=true`
- specify **working directory** (`dir`) and **timeout** (`timeout`)
- lightweight **CLI-first tool**, fully standalone

## 💡 Why

- **CLI-first and lightweight** — no editor or notebook required
- handles **dependencies, exports, and file outputs** out of the box
- focused on **Markdown workflows**, not just single snippets
- **extensible** to new languages and execution strategies
- ideal for **automation, literate programming, and reproducible docs**

## 👍 Support

All tools are completely free to use, with every feature fully unlocked and accessible.

If you find one or more of these tool helpful, please consider supporting its development with a donation.

Your contribution, no matter the amount, helps cover the time and effort dedicated to creating and maintaining these tools, ensuring they remain free and receive continuous improvements.

Every bit of support makes a meaningful difference and allows me to focus on building more tools that solve real-world challenges.

Thank you for your generosity and for being part of this journey!

[![Donate with PayPal](https://img.shields.io/badge/💸-Tip%20me%20on%20PayPal-0070ba?style=for-the-badge&logo=paypal&logoColor=white)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=FV575PVWGXZBY&source=url)

## 🛠️ How To Install

### Download the latest binaries from the [releases page](https://github.com/lucasepe/snipr/releases/latest):

- [macOS](https://github.com/lucasepe/snipr/releases/latest)
- [Windows](https://github.com/lucasepe/snipr/releases/latest)
- [Linux (arm64)](https://github.com/lucasepe/snipr/releases/latest)
- [Linux (amd64)](https://github.com/lucasepe/snipr/releases/latest)

### Using a Package Manager

» macOS » [Homebrew](https://brew.sh/)

```sh
brew tap lucasepe/cli-tools
brew install snipr
```

## 📂 Try it out

Run a demo Markdown file directly from GitHub:

```bash
curl -sL https://raw.githubusercontent.com/lucasepe/snipr/main/testdata/basic.md | snipr run
```

This demonstrates dependencies, exports, file output, and skipped blocks in action.
