 # 1. Edit your .go files in VS Code, save.

  # 2. Rebuild the image (note the trailing "." — build context is repo root)
  docker build -t lab04 -f docker/Dockerfile .

  # 3. Recreate containers with the new image
  docker-compose -f docker/docker-compose.yml up -d --force-recreate

  # 4. Test
  docker exec -it lab04-client bash
  #   then inside: cd /lab04/client && ./client_bin rpc put city London