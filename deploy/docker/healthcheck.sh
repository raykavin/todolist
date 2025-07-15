#!/bin/sh

# Configuration
HEALTH_ENDPOINT="${HEALTH_ENDPOINT:-https://localhost:3000/health}"
MAX_TIME="${HEALTH_CHECK_TIMEOUT:-10}"
APP_NAME="${APP_NAME:-khambalia}"
LOG_FILE="${LOG_DIR:-/var/log/${APP_NAME}}/healthcheck.log"

# Function to log messages
log_message() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE" 2>/dev/null || echo "$1"
}

# Function to parse JSON value by key (string only)
parse_json_value() {
    echo "$1" | grep -o "\"$2\":\"[^\"]*\"" | sed "s/\"$2\":\"\([^\"]*\)\"/\1/" | head -1
}

# Check if curl is available
if command -v curl >/dev/null 2>&1; then
    log_message "Starting healthcheck for ${HEALTH_ENDPOINT}"

    # Make the health check request
    RESPONSE=$(curl --insecure --fail --silent --max-time "$MAX_TIME" \
        --write-out "\n%{http_code}" \
        --header "Accept: application/json" \
        "$HEALTH_ENDPOINT" 2>/dev/null)

    # Extract HTTP status code
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    # Check if request was successful
    if [ "$HTTP_CODE" != "200" ] && [ "$HTTP_CODE" != "000" ]; then
        log_message "ERROR: Health endpoint returned HTTP $HTTP_CODE"
        exit 1
    fi

    # Check if we got a response body
    if [ -z "$BODY" ]; then
        log_message "ERROR: Empty response from health endpoint"
        exit 1
    fi

    # Parse expected JSON fields
    APP_VALUE=$(parse_json_value "$BODY" "app")
    STATUS_VALUE=$(parse_json_value "$BODY" "status")
    VERSION_VALUE=$(parse_json_value "$BODY" "version")

    # Validate fields
    if [ "$APP_VALUE" = "Todo List" ] && [ "$STATUS_VALUE" = "ok" ]; then
        log_message "Health check passed: app='$APP_VALUE', status='$STATUS_VALUE', version='$VERSION_VALUE'"
        exit 0
    else
        log_message "ERROR: Health check failed - unexpected values"
        log_message "Response: $BODY"
        exit 1
    fi

elif command -v wget >/dev/null 2>&1; then
    # Fallback to wget if curl is not available
    log_message "Using wget for healthcheck"

    RESPONSE=$(wget --quiet --no-check-certificate --timeout="$MAX_TIME" --tries=1 -O - "$HEALTH_ENDPOINT")

    APP_VALUE=$(echo "$RESPONSE" | grep -o '"app":"[^"]*"' | sed 's/"app":"\([^"]*\)"/\1/')
    STATUS_VALUE=$(echo "$RESPONSE" | grep -o '"status":"[^"]*"' | sed 's/"status":"\([^"]*\)"/\1/')

    if [ "$APP_VALUE" = "Todo List" ] && [ "$STATUS_VALUE" = "ok" ]; then
        log_message "Health check passed (wget)"
        exit 0
    else
        log_message "ERROR: Health check failed (wget)"
        log_message "Response: $RESPONSE"
        exit 1
    fi

else
    # If neither curl nor wget is available, check if process is running
    log_message "No HTTP client available, checking process"

    if pgrep -x "$APP_NAME" >/dev/null; then
        log_message "Process $APP_NAME is running"
        exit 0
    else
        log_message "ERROR: Process $APP_NAME not found"
        exit 1
    fi
fi
