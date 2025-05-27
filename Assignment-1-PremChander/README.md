# Gym management system
The GYM MANAGEMENT SYSTEM is to automate everything that happens in the gym. It is to aid and simplify the job all those who work for the gym, who train in the gym and who owns the gym. The database consists of daily goals and stats of achievements/progress of a member who trains in the gym, contact details and personal info of everyone, training programs that the gym offers, equipment etc.

## Who is it for
- Gym members
- Admin
- Receptionist
- Trainers

## Functionality
* **Admin** - Admin/Owner of the gym can add a receptionist, delete a receptionist, add a trainer, and delete a trainer, view details of everyone (Receptionist, Trainers, and gym members), add a new member to the gym, delete a member when one is leaving, add equipment, remove equipment.
* **Gym members** - Gym members can login into the system to view their daily goals like their workouts, goals, diet plan, etc. They enter their daily training details. They can edit their personal information like contact details, personal info etc. They can also view their trainers contact info. They can track their progress.
* **Trainers** - Trainer can design and add a workout plan, view all the details, progress of all those who train under him/her, evaluate their workouts every day; edit their own contact details and personal info; can view all the equipment details available in the gym.
 * **Receptionist** - Receptionist can view contact details of everyone, add a member, delete a member.

 PROG8860 Assignment 1 – CI/CD Pipeline with GitHub Actions

## 👨‍💻 Author
Prem Chander Jebastian
9015480

## 🚀 Technologies Used
- Python 3.9
- Flask
- MySQL (as a GitHub Actions service)
- Docker
- GitHub Actions
- pytest

## 📁 Folder Structure


CICD:


## ✅ CI Pipeline Includes:
1. **Build**: Install dependencies and setup Python
2. **Test**: Run 3 pytest-based test cases
3. **Containerize**: Docker image build
4. **Deploy**: Run container inside CI
5. **Monitor**: Print container status
6. **Teardown**: Stop container

## 🛠 How to Trigger
- Triggered on **Pull Request only**
- Create a branch (e.g., `assignment-1-feX-branch`)
- Push code
- Create PR and assign your instructor as a reviewer
