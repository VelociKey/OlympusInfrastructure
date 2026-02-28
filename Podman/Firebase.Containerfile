FROM google/cloud-sdk:latest

# Install GCP Emulators
RUN gcloud components install 
    pubsub-emulator 
    cloud-firestore-emulator 
    spanner-emulator 
    bigtable-emulator 
    --quiet

# Install Node.js and Firebase Tools
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash - && 
    apt-get install -y nodejs && 
    npm install -g firebase-tools

# Expose all common emulator ports
# 8080: Firestore
# 8085: Pub/Sub
# 9000: Realtime DB
# 9099: Auth
# 9199: Storage
# 4000: Hub
# 5001: Functions
# 9010: Spanner
# 9090: Bigtable
EXPOSE 8080 8085 9000 9099 9199 4000 5001 9010 9090

# Default entrypoint starts all requested emulators
# We will use environment variables to control which ones start
CMD ["bash"]
