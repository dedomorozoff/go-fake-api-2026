#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
EMAIL="test_$(date +%s)@example.com"
PASSWORD="Password123!"
NAME="TestUser"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Starting API tests at $API_URL..."

# Helper function to check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed. Please install jq to run this script.${NC}"
    exit 1
fi

# 1. Registration
echo -n "1. Testing Registration... "
REGISTER_RES=$(curl -s -X POST "$API_URL/registration" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"$NAME\",
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

if echo "$REGISTER_RES" | jq -e '.data.user' > /dev/null; then
    echo -e "${GREEN}SUCCESS${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "$REGISTER_RES"
    exit 1
fi

# 2. Authorization
echo -n "2. Testing Authorization... "
AUTH_RES=$(curl -s -X POST "$API_URL/authorization" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

TOKEN=$(echo "$AUTH_RES" | jq -r '.data.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo -e "${GREEN}SUCCESS${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "$AUTH_RES"
    exit 1
fi

# 3. Create Board
echo -n "3. Testing Create Board... "
CREATE_BOARD_RES=$(curl -s -X POST "$API_URL/boards" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"name\": \"My Test Board\",
        \"is_public\": true
    }")

BOARD_ID=$(echo "$CREATE_BOARD_RES" | jq -r '.data.id')

if [ "$BOARD_ID" != "null" ] && [ -n "$BOARD_ID" ]; then
    echo -e "${GREEN}SUCCESS (ID: $BOARD_ID)${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "$CREATE_BOARD_RES"
    exit 1
fi

# 4. List My Boards
echo -n "4. Testing List My Boards... "
LIST_RES=$(curl -s -X GET "$API_URL/boards" \
    -H "Authorization: Bearer $TOKEN")

if echo "$LIST_RES" | jq -e '.data | length > 0' > /dev/null; then
    echo -e "${GREEN}SUCCESS${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "$LIST_RES"
    exit 1
fi

# 5. Public Boards
echo -n "5. Testing Public Boards... "
PUBLIC_RES=$(curl -s -X GET "$API_URL/public-boards")

if echo "$PUBLIC_RES" | jq -e '.data' > /dev/null; then
    echo -e "${GREEN}SUCCESS${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "$PUBLIC_RES"
    exit 1
fi

echo -e "\n${GREEN}All tests passed successfully!${NC}"
