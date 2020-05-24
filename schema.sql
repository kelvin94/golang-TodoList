      --user
    CREATE TABLE  IF NOT EXISTS  task_management_user (
        id SERIAL UNIQUE,
        username varchar(100),
        password varchar(1000),
        email varchar(100)
    );
    INSERT INTO task_management_user VALUES(1,'suraj','suraj','sapatil@live.com');


    --category
    CREATE TABLE  IF NOT EXISTS category( 
        id SERIAL UNIQUE,
        name varchar(1000) not null, 
        task_management_user_id INT references task_management_user(id)
    );

    INSERT INTO "category" VALUES(1,'TaskApp',1);
    select * from category;

    --status
     
     
    CREATE TABLE  IF NOT EXISTS status(
        id SERIAL UNIQUE,
        status varchar(50) not null
    );
    INSERT INTO "status" VALUES(1,'COMPLETE');
    INSERT INTO "status" VALUES(2,'PENDING');
    INSERT INTO "status" VALUES(3,'DELETED');
     

    --task

     
     
    CREATE TABLE  IF NOT EXISTS task (
        id SERIAL UNIQUE,
        title varchar(100),
        content text,
        created_date timestamp,
        last_modified_at timestamp,
        finish_date timestamp,
        priority integer, 
        category_id INT references category(id), 
        task_status_id INT references status(id), 
        due_date timestamp, 
        task_management_user_id INT references task_management_user(id), 
        hide int
    );

    INSERT INTO "task" VALUES(1,'Publish on github','Publish the source of tasks and picsort on github','2015-11-12 15:30:59','2015-11-21 14:19:22','2015-11-17 17:02:18',3,1,1,NULL,1,0);
    INSERT INTO "task" VALUES(4,'gofmtall','The idea is to run gofmt -w file.go on every go file in the listing, *Edit turns out this is is difficult to do in golang **Edit barely 3 line bash script. ','2015-11-12 16:58:31','2015-11-14 10:42:14','2015-11-13 13:16:48',3,1,1,NULL,1,0);

    CREATE TABLE  IF NOT EXISTS comments(
        id SERIAL UNIQUE,
        content varchar(1000), 
        task_id INT references task(id), 
        created timestamp, 
        task_management_user_id INT references task_management_user(id)
    );

    CREATE TABLE  IF NOT EXISTS files(
        name varchar(1000) not null, 
        autoName varchar(255) not null, 
        task_management_user_id INT references task_management_user(id), 
        created_date timestamp);