#!/bin/bash

# Developer Sentiment Survey Script
# Dit script verzamelt developer sentiment data via een survey

set -e

# Configuration
SURVEY_DIR="monitoring/platform-value/surveys"
SURVEY_FILE="${SURVEY_DIR}/sentiment-$(date +%Y-%m).json"
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Create survey directory if it doesn't exist
mkdir -p "${SURVEY_DIR}"

# Survey Questions
echo "=== Developer Sentiment Survey ==="
echo ""
echo "Deze survey helpt ons het platform te verbeteren."
echo "Alle antwoorden zijn anoniem en worden alleen gebruikt voor platform verbetering."
echo ""

# Question 1: Overall Satisfaction (1-5)
echo "1. Hoe tevreden ben je over het platform? (1 = zeer ontevreden, 5 = zeer tevreden)"
read -p "Score (1-5): " SATISFACTION_SCORE

# Validate score
if ! [[ "$SATISFACTION_SCORE" =~ ^[1-5]$ ]]; then
    echo "❌ Ongeldige score. Gebruik een getal tussen 1 en 5."
    exit 1
fi

# Question 2: Ease of Use (1-5)
echo ""
echo "2. Hoe makkelijk vind je het platform te gebruiken? (1 = zeer moeilijk, 5 = zeer makkelijk)"
read -p "Score (1-5): " EASE_OF_USE

if ! [[ "$EASE_OF_USE" =~ ^[1-5]$ ]]; then
    echo "❌ Ongeldige score. Gebruik een getal tussen 1 en 5."
    exit 1
fi

# Question 3: Cognitive Load (1-10)
echo ""
echo "3. Hoe complex vind je het platform? (1 = zeer simpel, 10 = zeer complex)"
read -p "Score (1-10): " COGNITIVE_LOAD

if ! [[ "$COGNITIVE_LOAD" =~ ^[0-9]|10$ ]]; then
    echo "❌ Ongeldige score. Gebruik een getal tussen 1 en 10."
    exit 1
fi

# Question 4: Time to First Deployment
echo ""
echo "4. Hoe lang duurde het voordat je je eerste deployment deed? (in minuten)"
read -p "Tijd (minuten): " TIME_TO_FIRST_DEPLOYMENT

if ! [[ "$TIME_TO_FIRST_DEPLOYMENT" =~ ^[0-9]+$ ]]; then
    echo "❌ Ongeldige tijd. Gebruik een getal (minuten)."
    exit 1
fi

# Question 5: Self-Service Usage
echo ""
echo "5. Gebruik je self-service deployment? (ja/nee)"
read -p "Antwoord (ja/nee): " SELF_SERVICE_USAGE

if [[ ! "$SELF_SERVICE_USAGE" =~ ^(ja|nee)$ ]]; then
    echo "❌ Ongeldig antwoord. Gebruik 'ja' of 'nee'."
    exit 1
fi

# Question 6: Documentation Usage
echo ""
echo "6. Gebruik je de platform documentatie? (ja/nee)"
read -p "Antwoord (ja/nee): " DOCUMENTATION_USAGE

if [[ ! "$DOCUMENTATION_USAGE" =~ ^(ja|nee)$ ]]; then
    echo "❌ Ongeldig antwoord. Gebruik 'ja' of 'nee'."
    exit 1
fi

# Question 7: Feedback (optional)
echo ""
echo "7. Heb je feedback of suggesties voor verbetering? (optioneel)"
read -p "Feedback: " FEEDBACK

# Question 8: Most Valuable Feature
echo ""
echo "8. Wat is de meest waardevolle feature van het platform voor jou?"
read -p "Feature: " MOST_VALUABLE_FEATURE

# Question 9: Biggest Pain Point
echo ""
echo "9. Wat is je grootste pijnpunt met het platform?"
read -p "Pijnpunt: " BIGGEST_PAIN_POINT

# Create JSON response
RESPONSE=$(cat <<EOF
{
  "timestamp": "${TIMESTAMP}",
  "satisfaction_score": ${SATISFACTION_SCORE},
  "ease_of_use": ${EASE_OF_USE},
  "cognitive_load": ${COGNITIVE_LOAD},
  "time_to_first_deployment_minutes": ${TIME_TO_FIRST_DEPLOYMENT},
  "self_service_usage": "${SELF_SERVICE_USAGE}",
  "documentation_usage": "${DOCUMENTATION_USAGE}",
  "feedback": "${FEEDBACK}",
  "most_valuable_feature": "${MOST_VALUABLE_FEATURE}",
  "biggest_pain_point": "${BIGGEST_PAIN_POINT}"
}
EOF
)

# Append to survey file (create array if doesn't exist)
if [ ! -f "$SURVEY_FILE" ]; then
    echo "[]" > "$SURVEY_FILE"
fi

# Add response to JSON array
# Using jq if available, otherwise manual JSON manipulation
if command -v jq &> /dev/null; then
    jq ". += [${RESPONSE}]" "$SURVEY_FILE" > "${SURVEY_FILE}.tmp" && mv "${SURVEY_FILE}.tmp" "$SURVEY_FILE"
else
    # Fallback: append manually (simple but works)
    # Remove closing bracket, add comma, add response, add closing bracket
    sed -i '' '$ s/]$/,/' "$SURVEY_FILE" 2>/dev/null || sed -i '$ s/]$/,/' "$SURVEY_FILE"
    echo "$RESPONSE" >> "$SURVEY_FILE"
    echo "]" >> "$SURVEY_FILE"
fi

echo ""
echo "✅ Survey ingevuld! Bedankt voor je feedback."
echo "📊 Resultaten worden opgeslagen in: ${SURVEY_FILE}"
echo ""
echo "Deze data wordt gebruikt om het platform te verbeteren."

