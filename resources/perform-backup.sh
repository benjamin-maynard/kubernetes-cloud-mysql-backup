#/bin/sh

has_failed=false

if mysqldump -u $TARGET_DATABASE_USER -h $TARGET_DATABASE_HOST -p$TARGET_DATABASE_PASSWORD -P $TARGET_DATABASE_PORT $TARGET_DATABASE_NAME > /tmp/$AWS_BUCCKET_BACKUP_NAME
then
    
    echo -e "Database backup successfully completed for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."

    if aws s3 cp /tmp/$AWS_BUCCKET_BACKUP_NAME s3://$AWS_BUCKET_NAME$AWS_BUCKET_BACKUP_PATH$AWS_BUCCKET_BACKUP_NAME
    then
        echo -e "Database backup successfully uploaded for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."
    else
        echo -e "Database backup failed to upload for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."
        has_failed=true
    fi

else
    echo -e "Database backup FAILED for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."
    has_failed=true
fi



if [ "$has_failed" = true ]
then
    echo -e "aws-database-backup encountered 1 or more errors. Exiting with status code 1."
    exit 1
else
    exit 0
fi