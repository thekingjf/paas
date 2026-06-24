# Bugs & Debugging

## How to use this file
Log an entry only if it clears the bar:
- BUG: it taught me something or cost real time to find.
Keep entries 2–4 lines. Capture the WHY / the lesson, not the code.
When unsure, log it — under-logging (forgetting a good story) is the bigger risk.

---
## Template
### <short title of the bug>
- **Symptom:** <what you saw / what broke>
- **Cause:** <what it actually was>
- **Fix / lesson:** <what you did + the transferable takeaway>

### Package name
- **Symptom:** code wouldn't run
- **Cause:** wrong package name made it a library
- **Fix / lesson:** package main is what makes a runnable program

### Driver Imports
- **Symptom:** Errors when importing modernc
- **Cause:** the package is split into a driver and constants and to work you need both
- **Fix / lesson:** Added both imports learned the driver registration was a side
                    effect of init() so a missing driver fails at runtime if you 
                    only use the types and fails at compile if you only use
                    the registration