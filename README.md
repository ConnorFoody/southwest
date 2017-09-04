# southwest
Southwest assigns seats in the order you check in. This project allows you to set up a task to check into southwest automatically. It also tries to get around network latency by firing several requests before checkin time.  

usage is southwest <time> <full name> <record locator>

It schedules a series of "blasts" starting 250ms before your checkin time. Blasts occur every 100ms until 250ms after your checkin time. These are configurable in the code. 

The requests go through the mobile API southwest uses. Note that southwest changes their mobile API frequently so this may not be up to date. Each blast establishes a session leading up to the "blast period". During the blast period it tries to checkin, and confirm. A locking system manages the blasts to make sure that requests aren't fired after a "stage" is complete by any of the other checkin tasks. 
