#/bin/sh

if [ -v "$2"]
then
PAYLOAD="payload={\"channel\": \"$SLACK_CHANNEL\", \"username\": \"$SLACK_USERNAME\", \"text\": \"$1\", \"icon_emoji\": \":slack:\"}"
else
PAYLOAD="payload={\"channel\": \"$SLACK_CHANNEL\", \"username\": \"$SLACK_USERNAME\", \"text\": \"$1\`\`\`$2\`\`\`\", \"icon_emoji\": \":slack:\"}"
fi

curl -X POST --data-urlencode "$PAYLOAD" "$SLACK_WEBHOOK_URL"