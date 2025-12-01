# TODO
- [x] Add controller support
- [ ] Fix XBOX controller support
- [x] Add backgound and foreground layers
- [x] Add spikes and death (?) => spikes as a foreground element that the player collision detects
- [ ] Add sound
- [ ] Fix player clipping through the floor (tunneling problem, solved by implementing Continuous Collision Detection)
- [x] Fix player being transported to top of platform on lateral hit
- [x] Fix prop follow code
- [x] Add interact button (don't open door until interaction)
- [ ] Save/Load state logic
- [ ] Improve inventory follow logic
- [ ] Add key floating animation
- [ ] Add falling player animation
- [x] Add jump and landing player particles
- [ ] Add running and stopping player particles
- [ ] Add controller rumble on landing
- [ ] Add proper death event, do not immediately reset, play sound and visual feedback
- [ ] Sort out entities dependency chart, things like player.ProcessInput(..., activeGamepad int) are starting to smell
- [ ] Treat player.HandleCollisions with the respect it deserves, and map out all the actual cases
- [x] Fix jump VFX not starting at ground level
- [ ] Add landing VFX
- [ ] Fix level particles not being properly reset on game reset
- [ ] Make spikes not deadly on the sides (Create separate hitbox for damage)

# IDEAS
- [ ] Help mother animals get their cubs
- [ ] An object that lets you stand on clouds
