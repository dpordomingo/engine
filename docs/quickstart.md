# Quick Start

1. [Install Docker](#1-install-docker)
1. [Install source{d} Engine](#2-install-source-d-engine)
1. [Start source{d} Engine](#3-start-source-d-engine)
1. [Explore source{d} Engine](#4-explore-source-d-engine)
    * [Querying code](#querying-code)
        * [gitbase web interface](#query-web-interface)
        * [gitbase SQL CLI](#query-command-line-interface)
        * [database schema](#database-schema)
    * [Parsing code](#parsing-code)
        * [bblfsh web client](#parse-web-client)
        * [bblfsh CLI](#parse-command-line-interface)


## 1 Install Docker

Follow these instructions based on your OS:

### Docker for macOS

Follow the instructions at [Docker for macOS](https://store.docker.com/editions/community/docker-ce-desktop-mac).  
You may also use [Homebrew](https://brew.sh/):
```bash
$ brew cask install docker
```

### Docker for Linux

On Ubuntu, follow the instructions at [Docker for Ubuntu Linux](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce-1):
```bash
$ sudo apt-get update
$ sudo apt-get install docker-ce
```

On Arch Linux, follow the instructions at [Docker for Arch Linux](https://wiki.archlinux.org/index.php/Docker#Installation):
```bash
$ sudo pacman -S docker
```

### Docker for Windows

On Windows, follow the instructions at [Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows). Make sure to read the system requirements [here](https://docs.docker.com/docker-for-windows/install/).

Please note Docker Toolbox is not supported.


## 2 Install source{d} Engine

Download the **[latest release](https://github.com/src-d/engine/releases/latest)** for macOS (Darwin), Linux or Windows.

### Install Engine on macOS or Linux

Extract `srcd` binary from the release you downloaded, and move it into your bin folder to make it executable from any directory:

```bash
$ tar -xvf path/to/engine_REPLACEVERSION_linux_amd64.tar.gz
$ sudo mv path/to/engine_REPLACE_OS_amd64/srcd /usr/local/bin/
```

### Install Engine on Windows

*The support for Windows is currently experimental.*

*Please note that from now on we assume that the commands are executed in `powershell` and not in `cmd`. Running them in `cmd` is not guaranteed to work. Proper support may be added in future releases.*

Create a directory for srcd.exe and add it to your `$PATH`;  
_you need to run these commands in a powershell as administrator._
```powershell
mkdir 'C:\Program Files\srcd'
# Add the directory to the `%path%` to make it available from anywhere
setx /M PATH "$($env:path);C:\Program Files\srcd"
# Now open a new powershell to apply the changes
```

Extract the `srcd.exe` executable from the release you downloaded, and copy it into the directory you created in the previous step:
```powershell
mv \path\to\engine_windows_amd64\srcd.exe 'C:\Program Files\srcd'
```


## 3 Start source{d} Engine

Now it's time to initialize source{d} Engine and provide it with some repositories to analyze:

```bash
# Without a path Engine operates on the local directory,
# it works with nested or sub-directories.
$ srcd init

# You can also provide a path
$ srcd init <path>
```

**Note:**
Ensure that you initialize source{d} Engine every time you want to process a new repository.
Changes in the `init` working directory are not detected automatically.
Database indexes are not updated automatically when its contents change, so the index should be manually recreated.

**Note for macOS:**
Docker for Mac [requires enabling file sharing](https://docs.docker.com/docker-for-mac/troubleshoot/#volume-mounting-requires-file-sharing-for-any-project-directories-outside-of-users) for any path outside of `/Users`.

**Note for Windows:** Docker for Windows [requires shared drives](https://docs.docker.com/docker-for-windows/#shared-drives). Other than that, it's important to use a workdir that doesn't include any sub-directory whose access is not readable by the user running `srcd`. For example, using `C:\Users` as workdir, will most probably not work. For more details see [this issue](https://github.com/src-d/engine/issues/250).


## 4 Explore source{d} Engine

_For the full list of the commands supported by `srcd` and those
that have been planned, please read [commands.md](docs/commands.md)._

source{d} Engine provides interfaces to [query code repositories](#querying-code) and to [parse code](#parsing-code) into [Universal Abstract Syntax Trees](#babelfish-uast).

**Note:**
source{d} Engine will download and install Docker images on demand. Therefore, the first time you run some of these commands, they might take a bit of time to start up. Subsequent runs will be faster.


### Querying Code

#### Query Web Interface

To launch the [web client for the SQL interface](https://github.com/src-d/gitbase-web), run the following command and start executing queries:

```bash
# Launch the query web client
$ srcd web sql
```

This should open the web interface in your browser.
You can also access it directly at [http://localhost:8080](http://localhost:8080).

#### Query Command Line Interface

If you prefer to work within your terminal via command line, you can open a SQL REPL
that allows you to execute queries against your repositories by executing:

```bash
# Launch the query CLI REPL
$ srcd sql
```

If you want to run a query directly, you can also execute it as such:

```bash
# Run query via CLI
$ srcd sql "select count(*) as count from repositories;"
```

_You can find further sample queries in the [examples](examples/README.md) folder._

**Note:**
Engine's SQL supports a [UAST](#babelfish-uast) function that returns a Universal AST for the selected source text. UAST values are returned as binary blobs, and are best visualized in the `web sql` interface rather than the CLI where are seen as binary data.

#### Database Schema

If you want to know what the database schema look like, you can refer to the [diagram about gitbase entities and relations](https://docs.sourced.tech/gitbase/using-gitbase/schema).

You can also query the tables that are available in source{d} Engine.

```bash
$ srcd sql "SHOW tables;"
+--------------+
|    TABLE     |
+--------------+
| blobs        |
| commit_blobs |
| commit_files |
| commit_trees |
| commits      |
| files        |
| ref_commits  |
| refs         |
| remotes      |
| repositories |
| tree_entries |
+--------------+
```

```bash
$ srcd sql "DESCRIBE TABLE commits;"
+---------------------+-----------+
|        NAME         |   TYPE    |
+---------------------+-----------+
| repository_id       | TEXT      |
| commit_hash         | TEXT      |
| commit_author_name  | TEXT      |
| commit_author_email | TEXT      |
| commit_author_when  | TIMESTAMP |
| committer_name      | TEXT      |
| committer_email     | TEXT      |
| committer_when      | TIMESTAMP |
| commit_message      | TEXT      |
| tree_hash           | TEXT      |
| commit_parents      | JSON      |
+---------------------+-----------+
```

### Parsing Code

Sometimes you may want to parse files directly as [UASTs](#babelfish-uast).

To see which languages are available, check the table of [supported languages](#babelfish-uast).

#### Parse Web Client

If you want a playground to see examples of the UAST, or run your own, you can launch the [parse web client](https://github.com/bblfsh/web).


```bash
# Launch the parse web client
$ srcd web parse
```

This should open the [web interface](https://github.com/bblfsh/web) in your browser.
You can also access it directly at [http://localhost:8081](http://localhost:8081).

#### Parse Command Line Interface

Alternatively, you can also start parsing files on the command line:

```bash
# Parse file via CLI
$ srcd parse uast /path/to/file.java
```

To parse a file specifying the programming language:

```bash
$ srcd parse uast --lang=LANGUAGE /path/to/file
```

To see the installed language drivers:

```bash
$ srcd parse drivers list
```

<!--
### 5. Next steps

You can now run source{d} Engine, choose what you would like to do next:

- [**Analyze your git repositories**](#)
- [**Understand how your code has evolved**](#)
- [**Write your own static analysis rules**](#)
- [**Build a data pipeline for MLonCode**](#)
-->