## student-schedule-app
A lightweight web-based shift scheduling system built with Go, allowing users to create weekly schedules and managers to review and approve them.  This project demonstrates full-stack functionality using Go’s net/http, HTML templates, and simple file-based persistence (no database required).

## Shift Scheduler Web App

A lightweight web-based shift scheduling system built with Go, allowing users to create weekly schedules and managers to review and approve them.

This project demonstrates full-stack functionality using Go’s net/http, HTML templates, and simple file-based persistence (no database required).

## Features
- Schedule Input
- Weekly view (Monday–Friday, 8 AM – 6 PM)
- Click-and-drag interface for selecting shifts
- Time divided into 10-minute increments
- Real-time tracking of:
  - Daily hours
  - Weekly total hours

## Constraints
Enforced on both frontend and backend logic:

- Minimum shift per day: 3 hours
- Maximum shift per day: 9 hours
- Weekly total: 20–40 hours

Invalid selections are visually highlighted and prevented from submission.

## Submission & Approval Flow
Users submit schedules for review
Status is displayed as:
- Pending
- Approved
- Denied

# Manager role can:
- Approve schedules
- Deny schedules

## Authentication (Basic)
Student login

  `Username: _________*`
  
  `Password: student12`
  
  _The students username can be anything teh first time, as it will create for them a personal save file_

Manager login

  `Username: manager`
  
  `Password: admin12`

Note: This is a simple demo authentication system (no hashing or sessions).

## Project Structure
shift-scheduler/
├── main.go                 # Main Go application

├── go.mod                  # Go module file

├── go.sum                  # Optional, dependency checksum file

├── static/                 # All CSS, JS, images

│   ├── app.js

│   └── style.css

├── templates/              # HTML templates

│   ├── login.html

│   ├── menu.html

│   ├── calendar.html

│   └── schedule.html

├── saveFiles/              # User schedule data (text files)

│   ├── cameronS.txt

│   ├── lazyLary.txt

│   ├── pendingPenny.txt

│   └── sallySmith.txt

├── handlers/               # Optional: Go code split by functionality

│   ├── auth.go

│   ├── schedule.go

│   └── manager.go

├── utils/                  # Optional: helper functions
│   └── fileops.go

└── README.md               # Documentation
    
## Setup & Installation
#1. Install Go

Make sure you have Go installed:

`go version`

If not, download from Go's [Official Website](https://go.dev/dl/)

#2. Clone the Repository

`git clone <your-repo-url>
cd shift-scheduler`

#3. Run the Application

`go run main.go`

#4. Open in Browser

http://localhost:8080

## How to Use

# Employee Flow
  Login with
  
  `Username: _________`
  
  `Password: student12`
  
  - Click Employee Time Sheet
  - Select time blocks by dragging
  - Ensure constraints are satisfied
  - Click Save
  - See status: Pending Approval

# Manager Flow
  Login as:
  
  `Username: manager`
  
  `Password: admin12`
  
  Click Schedule Status
  Select a user schedule
  View in read-only mode
  Select:
  
  `✅ Approve`
  
  `❌ Deny`
  
## Data Storage (Important)

This app uses plain text files instead of a database.

Each user has a file in:

/saveFiles/<username>.txt

## File Format
Pending
Mon-1-0
Mon-1-1
Tue-3-2
...

# Breakdown:
First line → Status (Pending, Approved, Denied)

Remaining lines → Selected time slices:

Day-HourIndex-SliceIndex
Example:
Wed-4-2
Wednesday
Hour index (from 8 AM start)
10-minute slice within that hour

## Limitations / Notes

No database (file-based only)
No real authentication/session management
Global variables used for user state (not production-safe)
No concurrency protection on file writes
Minimal validation on server (can be extended)
