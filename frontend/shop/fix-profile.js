// fix-profile.js
import fs from 'fs';
import path from 'path';

// Fix the profile.vue file
const fixProfileFile = () => {
    try {
        const profilePath = '/Users/dercyc/go/src/pro/mall-go/frontend/shop/src/views/user/profile.vue';

        // Read the file content
        let content = fs.readFileSync(profilePath, 'utf8');

        // Fix import statement
        content = content.replace(
            /import { getUserInfo, updateUserInfo, updatePassword } from '@\/api\/user';/g,
            "import { getUserInfo, updateUserInfo, updatePassword as updatePasswordApi } from '@/api/user';"
        );

        // Fix function definition
        content = content.replace(
            /\/\/ 修改密码\nconst updatePassword = async \(\) => {/g,
            "// 修改密码\nconst handleUpdatePassword = async () => {"
        );

        // Fix function call
        content = content.replace(
            /@click="updatePassword"/g,
            '@click="handleUpdatePassword"'
        );

        // Fix API call
        content = content.replace(
            /await updatePassword\({/g,
            "await updatePasswordApi({"
        );

        // Write the changes back to the file
        fs.writeFileSync(profilePath, content);
        console.log('Fixed: profile.vue');

    } catch (error) {
        console.error('Error fixing profile.vue:', error);
    }
};

// Run the fix
fixProfileFile();
