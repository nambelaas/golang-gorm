create table sample(
    id varchar(10) primary key,
    name varchar(50) not null,
)engine=InnoDB default charset=utf8mb4;

select * from sample

create table users(
    id int AUTO_INCREMENT,
    name varchar(50) not null,
    password varchar(50) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine=InnoDB default charset=utf8mb4;

alter table users
    rename column name to first_name;
alter table users
    add column middle_name varchar(50) after first_name,;
alter table users
    add column last_name varchar(50) after middle_name;

create table user_logs(
    id int AUTO_INCREMENT,
    user_id varchar(100) not null,
    action varchar(100) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine=InnoDB default charset=utf8mb4;

create table todos(
    id bigint not null AUTO_INCREMENT,
    user_id int not null,
    title varchar(100) not null,
    description text,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
) engine=InnoDB default charset=utf8mb4;

create table wallets(
    id varchar(100) not null,
    user_id varchar(100) not null,
    balance bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (user_id) references users(id) on delete cascade
) engine=InnoDB default charset=utf8mb4;

create table addresses(
    id bigint not null auto_increment,
    user_id int not null,
    address varchar(100) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
    foreign key (user_id) references users(id) on delete cascade
) engine=InnoDB default charset=utf8mb4;

create table products(
    id bigint not null auto_increment,
    name varchar(100) not null,
    price bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine=InnoDB default charset=utf8mb4;

create table user_like_product(
    user_id int not null,
    product_id int not null,
    primary key (user_id, product_id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (product_id) references products(id) on delete cascade
) engine=InnoDB default charset=utf8mb4;