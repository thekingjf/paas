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

### Go Docker Modules
- **Symptom:** go get github.com/docker/docker/client would return module declares
               its path as github.com/moby/moby/client but was required as 
               github.com/docker/docker/client.
- **Cause:** docker/docker/client used to just be a package in docker/docker but
             when they migrated it split the client into moby/moby/client
- **Fix / lesson:** pinned the parent docker/docker@v28.5.2+incompatible which is
                    old and unused will definitely modernize soon and also ran
                    go mod tidy. Learned a bug like declares its path as X but 
                    required as Y usually means a package got renamed/moved and to
                    fix it use the new path or pin a parent version that still
                    uses it
                    
### 
- **Symptom:** POST /apps/blog/deploy build failed with 500 / build failed. Server 
               log showed the real error: Error response from daemon: lsetxattr 
               /Dockerfile: xattr "com.apple.provenance": operation not supported
- **Cause:** The daemon failed to unpack apple specific attributes since they aren't
             recognized on linux systems
- **Fix / lesson:** rebuilt the tar without the attributes. Lesson learned was 
                    different tarring tools create different tars when given the
                    same content