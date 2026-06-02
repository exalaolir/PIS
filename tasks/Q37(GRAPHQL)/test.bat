// для теста не забыть поменять свои данные
@echo off
chcp 65001 >nul
set u=http://localhost:3000/graphql
set h=-H "Content-Type: application/json"

echo 1. Query celebrities (все)
curl -s -X POST %u% %h% -d "{\"query\":\"{ celebrities { id fullname nationality reqPhotoPath } }\"}"
echo.

echo 2. Mutation createCelebrity
curl -s -X POST %u% %h% -d "{\"query\":\"mutation { createCelebrity(input: { id: 10, fullname: \\\"Test Star\\\", nationality: \\\"Test\\\", reqPhotoPath: \\\"/photos/test.jpg\\\" }) { id fullname } }\"}"
echo.

echo 3. Query celebrity(id: 10)
curl -s -X POST %u% %h% -d "{\"query\":\"{ celebrity(id: 10) { id fullname nationality reqPhotoPath } }\"}"
echo.

echo 4. Mutation updateCelebrity
curl -s -X POST %u% %h% -d "{\"query\":\"mutation { updateCelebrity(id: 10, input: { fullname: \\\"Updated Star\\\", nationality: \\\"Updated\\\", reqPhotoPath: \\\"/photos/updated.jpg\\\" }) { id fullname } }\"}"
echo.

echo 5. Mutation deleteCelebrity
curl -s -X POST %u% %h% -d "{\"query\":\"mutation { deleteCelebrity(id: 10) }\"}"
echo.

pause
