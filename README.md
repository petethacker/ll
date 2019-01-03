a worthwhile README.md is on the way.

ll is aimed at windows only because i find i type `ll` on windows a million times per day when i'm swappig between linux and windows systems. There is no need for this on linux, so its not being tested or developed with linux in mind.

Yes i could have just used some kind of syslink to point to `dir` or just wrapped `dir` into a subprocess within GO, but that not as fun as just making something!

No arguments are supported yet, only arg 1 which is the path to search. If no path is specified then it uses the current working directory.