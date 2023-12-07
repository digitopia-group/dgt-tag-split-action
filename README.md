## Usage

This action is meant for splitting a specifically formatted tag into its parts.  
For example: "client_v0.0.3"  
The action will return 2 variables: "tag" will be set to "client" and "versionnr" will be set to "0.0.3".  

This is handy for monorepos that need to build some specific subprogram.  

Maybe there is a better way of doing this, if you know one, please let me know.


### Example workflow

```yaml
name: Tag splitting and getting of parts
on: 
  push:
    tags:
      - '*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: Run action

      # Put your action repo here
      uses: houseofdigital/dgt-tag-split-action

````
### Inputs

| Input                | Description                                                                     |
|----------------------|---------------------------------------------------------------------------------|
| `refname`            | The github.ref_name. This is filled in by default, you can change it, but why?  |

### Outputs

| Output                                               | Description                                   |
|------------------------------------------------------|-----------------------------------------------|
| `tag`                                                | the first part of the tag                     |
| `versionnr`                                          | the version that comes out of it.             |

## Examples

```yaml
steps:
- uses: actions/checkout@master
- name: Run action
  id: myaction

  uses: houseofdigital/dgt-tag-split-action

- name: Check outputs
    run: |
    echo "Outputs - ${{ steps.myaction.outputs.tag }}"
    echo "Outputs - ${{ steps.myaction.outputs.versionnr }}"

```
