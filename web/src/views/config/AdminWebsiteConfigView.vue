<script setup lang="ts">
import { ConfigAPI, type ConfigDTO } from '@/api/config.api';
import { UserAPI, UserRole } from '@/api/user.api';
import router from '@/router';
import type { AxiosError } from 'axios';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const configDTO = ref<ConfigDTO | undefined>(undefined);
const error = ref<string | undefined>(undefined);
const route = useRoute();

const getConfig = () => {
    error.value = undefined;

    ConfigAPI.GetConfig().then((res) => {
        configDTO.value = res.data.data!;
    }).catch((err) => {
        error.value = "Failed to load config from the backend";
        console.error(err);
    });
}

onMounted(() => {
    UserAPI.AuthorizeUser({required_roles: [UserRole.Admin]}).then(() => {
        getConfig();
    }).catch((err: AxiosError) => {
        console.error(err);
        router.push({ path: '/login', query: { redirect: route.fullPath } });
    });
});

const onCheckPostingEnabled = () => {
    if (!configDTO.value) {
        return;
    }
    ConfigAPI.SetPostingEnabled(configDTO.value.posting_enabled).then((res) => {
        configDTO.value = res.data.data!;
    }).catch((err) => {
        error.value = "Failed to update configuration";
        console.error(err);
    });
}

</script>

<template>
    <h1>Website configuration</h1>

    <div v-if="configDTO">
        <label for="posting_enabled">Posting enabled:</label>
        <input id="posting_enabled" type="checkbox" v-model="configDTO.posting_enabled" @change="onCheckPostingEnabled" />
        <br/>
        <small>
            If checked, users can create threads and posts.
            If unchecked, posting is disabled until further notice.
        </small>
    </div>

    <div class="error" v-if="error">{{ error }}</div>
</template>

<style lang="css" scoped>
.error {
    color: var(--user-error-color);
}
</style>