@echo off
REM Migration script for Windows

SET DIRECTION=%1
IF "%DIRECTION%"=="" SET DIRECTION=up

IF "%DATABASE_URL%"=="" (
    echo Error: DATABASE_URL environment variable is not set
    exit /b 1
)

echo Running migrations: %DIRECTION%

SET SCRIPT_DIR=%~dp0
SET MIGRATIONS_DIR=%SCRIPT_DIR%..\migrations

IF "%DIRECTION%"=="up" (
    FOR %%F IN (%MIGRATIONS_DIR%\*.up.sql) DO (
        echo Running migration: %%~nxF
        psql "%DATABASE_URL%" -f "%%F"
    )
    echo Migrations completed successfully
) ELSE IF "%DIRECTION%"=="down" (
    FOR /L %%G IN (5,-1,1) DO (
        IF EXIST "%MIGRATIONS_DIR%\00%%G_*.down.sql" (
            FOR %%F IN (%MIGRATIONS_DIR%\00%%G_*.down.sql) DO (
                echo Running migration: %%~nxF
                psql "%DATABASE_URL%" -f "%%F"
            )
        )
    )
    echo Rollback completed successfully
) ELSE (
    echo Error: Invalid direction. Use 'up' or 'down'
    exit /b 1
)
