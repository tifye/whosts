whosts is a small CLI tool written in Go to edit Windows hosts files. The hosts file is a simple text file that acts as a local DNS lookup local to the machine its on. 

```text
Usage
  whosts [command]

Available Commands:
  add        Add an entry
  completion Generate the autocompletion script for the specified shell
  dump       Dumps file contents to stdout
  help       Help about any command
  list       List all entries
  open       Opens the hosts file in notepad
  remove     Remove entries matching passed filters. Filters are stacked
    --comment           Remove entries with matching comment
    --dry               Dry run command and print out which entries would have been removed
    --duplicates-only   Remove entry duplicates that match passed filters. If no filters are passed then remove any duplicate.
    --host              Remove entries with matching host name
    --ip                Remove entries with matching IP
    --no-comment        Remove entries without comments

Use "whosts [command] --help" for more information about a command.
```
