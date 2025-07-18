<div align="center">
  <img src="assets/image.png" alt="alt text" />
</div>

# idea
A very simple idea log that can sync up with notion. Built using golang. 

## Why does this exist? 
I get a crap ton of ideas every week, and I dont have the energy (and RAM - its notion we're talking about) to keep a billion notion tabs open. I also love the command line and wanted to build a CLI, so it really was a perfect blend of things to give me what I need. 

It's a very basic and simple project, and I really intended for it to be as simple as possible for other poeple to use it or work with it (if you're insane enough). 

## What did I use? 
1. bolt as my DB for local
2. bubbletea for the UI (some components still use survey)
3. notion integrations 
4. cobra for the CLI

## TODO: 
1. Have sha256 up every API keys ass to protect them from "malicious people (aka me)"
2. Refactor the whole thing considering go's best practices in mind
3. UNIT TESTING EVERYTHING AND EVERYWHERE
4. Testing for speed (yes, my 5 ideas per week should be uploaded fast, isnt speed nice?)
5. Integrating bubble tea and removing survey