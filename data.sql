CREATE TABLE IF NOT EXISTS toggle_features (
    "id" UUID DEFAULT uuid_generate_v1() PRIMARY KEY,
    "feature" character varying(255) not null,
    "is_enabled" bool default false,
    "created" timestamp default now(),
    "updated" timestamp
);

CREATE TRIGGER set_timestamp_utc
BEFORE UPDATE ON toggle_features
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp_utc();

INSERT INTO toggle_features
    (id, feature, is_enabled)
VALUES
    (uuid_generate_v4(), 'FEATURE_A', true),
    (uuid_generate_v4(), 'FEATURE_B', false),
    (uuid_generate_v4(), 'FEATURE_C', true);


-- trigger
CREATE OR REPLACE FUNCTION trigger_toggle_feature_change()
RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
         IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
        PERFORM pg_notify('toggle_feature',notification::text);
        
        RETURN NULL; 
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER toggle_feature_change
AFTER INSERT OR UPDATE OR DELETE ON toggle_features
    FOR EACH ROW EXECUTE PROCEDURE trigger_toggle_feature_change();


    
