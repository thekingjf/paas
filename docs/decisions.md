# Design Decisions

## How to use this file
Log an entry only if it clears the bar:
- DECISION: a competent engineer could have reasonably chosen otherwise.
Keep entries 2–4 lines. Capture the WHY / the lesson, not the code.
When unsure, log it — under-logging (forgetting a good story) is the bigger risk.

---
## Template
### <short title of the decision>
- **Fork:** <the choice you faced, one line>
- **Chose:** <what you picked / what you rejected>
- **Why:** <the reasoning — the defensible part>

### Driver Choice
- **Fork:** Had to decide between modernc vs mattn for my SQLite drivers for GO
- **Chose:** I chose modernc
- **Why:** Modernc simplifies deployment and allows future me to have an easier time
           when it comes up and as a personal project I don't need the best I need
           something that works consistently and easily without becoming it's own
           issue
    
### Server Struct Design
- **Fork:** Handlers needed access to the DB and 3 valid options:
            - Package level global DB var
            - pass db into every handler
            - Hang it onto a struct and with handler methods
- **Chose:** I chose to hang it on a struct with handler methods
- **Why:** A struct scales cleanly as dependcies increase and make handlers testable

### Parameterized Queries
- **Fork:** Fmt.sprintf vs ? placeholders for the sql string
- **Chose:** I chose to move forward with ? placeholders
- **Why:** I chose to do the ? placeholders for security because it passes the
           value over to the driver separately meaning it's never processed as SQL,
           preventing SQL injection that sprintf would have

### No Nulls Default Fill
- **Fork:** Inserting basic placeholder values vs leaving them NULL
- **Chose:** I chose to insert basic placeholder values("", 0, "created")
- **Why:** I chose to insert placeholders because a created row in this databse 
           implies existence so the placeholder values are geninune and not covering  
           for missing data.

### Errors.Is vs Errors.As
- **Fork:** Recognizing which to use
- **Chose:** I chose to use Errors is
- **Why:** When I need to check an error I know it's between if I need to match
           it(Is) or if I need to extract it(As)

### Docker Migration
- **Fork:** stay on docker/docker@v28(old) or migrate moby/moby/client
- **Chose:** I chose to migrate to the newer docker version instead of staying
             on the old frozen version
- **Why:** I chose to migrate to the newer version because it's reflective of the
           current state of docker which is the path docker supports moving forward
           and since it's early on I won't need to make deep changes and can avoid 
           potential bugs that the old version help that were fixed in the newer 
           version before i ever run into them. Although the trade off is mostly
           coverage since most tutorials/articles will be based off of the older
           version

### Tar Contents of Context Directory
- **Fork:** Tar content of context directory or the directory itself
- **Chose:** I choose to tar the content of the context directory
- **Why:** It matches how docker build treats context, so the docker paths
           resolve but in exchange you must point at whose contents are the root
           and if you give it at the wrong level you silently get wrong paths

### File in Memory
- **Fork:** File in memory(os.ReadFile) vs Streamed(io.Copy)
- **Chose:** I choose to store the files in memory
- **Why:** Contexts are small and simpler but if given a huge file it would all
           be loaded into the RAM but context are source files, so in practice
           it won't ever be too bad

### Tarring Locaiton
- **Fork:** Tarring on client vs server
- **Chose:** I chose to tar the files on the clients side before sending it to the
             server
- **Why:** I chose to tar the files on the clients side because it gives the server
           a tar file it can immediately work with instead of running around and
           tarring the files itself but in exchange the server now trust whatever
           the client sends over and this would make it incompatible on a browser

### Route Separation
- **Fork:** Edit existing create function for deploy or giving deploy it's own function
- **Chose:** I chose to give deploy it's own function
- **Why:** i chose to give deploy it's own function because someone can create
           an app but not deploy it and also creating a app tends to only happen
           once while deploying happen many times

### Stream Error
- **Fork:** Use stdlib decoder or moby's jsonmessage helper
- **Chose:** I chose to use the stdlib decoder
- **Why:** I chose that because I experimented with using moby's jsonmessage helper
           and it broke the project by brining in a incompatible version 

### Traffic routing
- **Fork:** middleware in front of chi vs host based routing inside chi
- **Chose:** I chose middleware in front of chi
- **Why:** I chose middleware in front of chi because it gave a concrete line
           of to split traffic to apps and commands about apps. With that i did
           accept that proxied request bypass all the chi middleware

### URL Shape
- **Fork:** Fixed vs arbitrary
- **Chose:** I choose fixed names (<name>.localhost) against arbitrary
- **Why:** I choose fixed names over arbitrary because it allows me to easily
           handle all inputs with the only big cost being no custom names

### App source name
- **Fork:** basename of the deploy dir (filepath.Base(filepath.Clean(dir))) over
            an explicit deploy <name> <dir> arg.
- **Chose:** I chose to do the basename of the deploy dir
- **Why:** Matches Vercel/Heroku convention over configuration. With this choice
           name is coupled to the folder, deploying one dir under two names means
           renaming it
- **Note:** spent ~10 min deliberating a two-line reversible choice. Calibration 
            lesson: reversible one-file decisions get minutes; deliberation 
            budget should track cost-of-being-wrong 