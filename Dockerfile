# Set the base image
FROM alpine:3.11

RUN apk -v --update add \
        python3 \
        py-pip \
        groff \
        less \
        mailcap \
        mysql-client \
        curl \
        py-crcmod \
        bash \
        libc6-compat \
        gnupg \
        coreutils \        
        gzip \
        && \
    pip3 install --upgrade awscli s3cmd python-magic && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/*

# Set Default Environment Variables
ENV TARGET_DATABASE_PORT=3306
ENV SLACK_ENABLED=false
ENV SLACK_USERNAME=kubernetes-s3-mysql-backup
ENV CLOUD_SDK_VERSION=285.0.1
ENV BACKUP_PROVIDER=aws

# Set Google Cloud SDK Path
ENV PATH /google-cloud-sdk/bin:$PATH

# Install Google Cloud SDK
RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud --version

# Copy Slack Alert script and make executable
COPY resources/slack-alert.sh /
RUN chmod +x /slack-alert.sh

# Copy backup script and execute
COPY resources/perform-backup.sh /
RUN chmod +x /perform-backup.sh
CMD ["sh", "/perform-backup.sh"]