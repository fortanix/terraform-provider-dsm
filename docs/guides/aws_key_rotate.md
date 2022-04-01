## There are 3 ways to achieve rotation

### 1. Rotate local key and copy to AWS group

* key material in DSM.
* There will not be key link to old Local key.
* Aliases will automatically be moved to new key.

### 2. Rotate with 'DSM' option
* key material in DSM.
* There will be key link to old Local key.
* Aliases will automatically be moved to new key.


### 3. Rotate with 'AWS' option 
* key material not in DSM
* This does not support aliases.