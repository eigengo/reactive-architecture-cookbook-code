create keyspace eas with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

create table eas.session_sensor_values(
  session_id UUID,
  sequence int,
  sensor_values list<float>,
  valid boolean,
  primary key (session_id, sequence)
) with clustering order by (sequence asc);

create table eas.session_labels(
  session_id UUID,
  sequence int,
  start_time float,
  duration float,
  label text,
  automatic boolean,
  valid boolean,
  primary key (session_id, automatic, sequence)
);

create table eas.session_sensors(
  session_id UUID,
  sequence int,
  location int,
  data_types set<int>,
  valid boolean,
  primary key (session_id, sequence)
) with clustering order by (sequence asc);
