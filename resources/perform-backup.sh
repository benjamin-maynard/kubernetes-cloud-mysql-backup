#/bin/sh


# Set the has_failed variable to false. This will change if any of the subsequent database backups/uploads fail.
has_failed=false


# Perform the database backup. Put the output to a variable. If successful upload the backup to S3, if unsuccessful print an entry to the console and the log, and set has_failed to true.
if sqloutput=$(mysqldump -u $TARGET_DATABASE_USER -h $TARGET_DATABASE_HOST -p$TARGET_DATABASE_PASSWORD -P $TARGET_DATABASE_PORT $TARGET_DATABASE_NAME 2>&1 > /tmp/$AWS_BUCCKET_BACKUP_NAME)
then
    
    echo -e "Database backup successfully completed for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."

    # Perform the upload to S3. Put the output to a variable. If successful, print an entry to the console and the log. If unsuccessful, set has_failed to true and print an entry to the console and the log
    if awsoutput=$(aws s3 cp /tmp/$AWS_BUCCKET_BACKUP_NAME s3://$AWS_BUCKET_NAME$AWS_BUCKET_BACKUP_PATH$AWS_BUCCKET_BACKUP_NAME 2>&1)
    then
        echo -e "Database backup successfully uploaded for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S')."
    else
        echo -e "Database backup failed to upload for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S'). Error: $awsoutput" | tee -a /tmp/aws-database-backup.log
        has_failed=true
    fi

else
    echo -e "Database backup FAILED for $TARGET_DATABASE_NAME at $(date +'%d-%m-%Y %H:%M:%S'). Error: $sqloutput" | tee -a /tmp/aws-database-backup.log
    has_failed=true
fi


# Put the contents of the database backup logs into a variable
logcontents=`cat /tmp/aws-database-backup.log`


# Check if any of the backups have failed. If so, exit with a status of 1. Otherwise exit cleanly with a status of 0.
if [ "$has_failed" = true ]
then

    # If Slack alerts are enabled, send a notification alongside a log of what failed
    if [ "$SLACK_ENABLED" = true ]
    then
        /slack-alert.sh "One or more backups on database host $TARGET_DATABASE_HOST failed. The error details are included below:" "$logcontents"
    fi

    echo -e "aws-database-backup encountered 1 or more errors. Exiting with status code 1."
    exit 1

else

    # If Slack alerts are enabled, send a notification that all database backups were successful
    if [ "$SLACK_ENABLED" = true ]
    then
        /slack-alert.sh "All database backups successfully completed on database host $TARGET_DATABASE_HOST."
    fi

    exit 0
    
fi