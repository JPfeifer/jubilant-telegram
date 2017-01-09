Add simple Go script that lets you add Time entry's to Redmine via the command line. 
Usage: 
-activity_id string redmine activity_id 
-comment string comment for redmine i
-hours string hours spend on tasks 
-password string login password 
-project_id string redmine project_id 
-url string Server IP (default "0.0.0.0") 
-user string login user 
-user_id string redmine user_id 
or create a .redmine file in your home folder and fill it like this: 
url=0.0.0.0/redmine 
user=admin 
password=password 
project_id=1 
user_id=1 
activity_id=8 
The hours Attribute is always requiert.
