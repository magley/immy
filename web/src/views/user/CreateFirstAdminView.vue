<script setup lang="ts">
    import type { ApiErrorInfo, ApiResponse } from '@/api/http';
    import { UserAPI, type CreateFirstAdminDTO, type UserDTO } from '@/api/user.api';
    import type { AxiosError } from 'axios';
    import { ref } from 'vue';
    import { useRouter } from 'vue-router';

    const dto = ref<CreateFirstAdminDTO>({
        username: '',
        password: ''
    });
    const error = ref<string | undefined>(undefined);
    const router = useRouter();

    const onSubmitCreateFirstAdmin = () => {
        error.value = undefined;
        UserAPI.CreateFirstAdmin(dto.value).then((res) => {
            router.push("/login");
        }).catch((err: AxiosError<ApiResponse<UserDTO>>) => {
            error.value = err.response?.data.error?.message;
        });
    }
</script>

<template>
    <h1>Create the first admin</h1>

    <p id="info">
        Create the very first admin in the system.
        <br/><br/>
        Note that this form is accessible by anyone, so make sure the owner of this
        website is the one who creates the admin.<br/>If you are a regular user, close
        this page and return back <a href="/">home</a>.
        <br/><br/>
        Only the very first admin can be created here. Further attempts will fail,
        so make sure you remember the password.
    </p>

    <form @submit.prevent="onSubmitCreateFirstAdmin">
        <h1>Create First Admin</h1>
        
        <input type=text placeholder="Username" required v-model="dto.username"/>
        <br/>
        <input type=password placeholder="Password" required v-model="dto.password"/>
        <br/>
        <br/>
        <button type=submit>Submit</button>
        
        <template v-if="error">
            <div/>
            <span class="error">{{error}}</span>
        </template>
    </form>
</template>

<style scoped>
    h1 {
        text-align: center;
    }
    
    p#info {
        text-align: center;
        width: 30%;
        margin: auto;
    }

    form {
		width: 40%;
		margin: 2em auto;
		border: 1px solid black;
		background-color: var(--background-color-darker);
		padding: 1em;
		text-align: center;
	}

	h1 {
		color: var(--banner-title-color);
	}
	.error {
		color: var(--user-error-color);
	}
</style>