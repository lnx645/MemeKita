pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}
dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}
rootProject.name = "MemeKita"
include(":app")
include(":feature:feed")
include(":feature:duel")
include(":feature:challenge")
include(":feature:coins")
include(":feature:profile")
