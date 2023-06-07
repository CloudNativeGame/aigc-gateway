# AIGC-Gateway

This project aims to address the resource management issues of AIGC instances by providing an AIGC serverless gateway based on the auto-scaling feature of cloud-native architecture.

The gateway has the following features:

- User management. Each user has their own AIGC instance, and the gateway will maintain the mapping between the user and the instance.
- User-level resource management. AIGC computing instances are created and destroyed based on the user's login/offline status, while preserving user data.

## Dashboard of AIGC-Gateway

Click on Start to use.

![Initial page](./docs/images/dashboard-login.png)

Login or Register.

![user login page](./docs/images/user-login.png)

The new user does not own an AIGC instance, can choose to install the required instance.

![instance uninstalled page](./docs/images/dashboard-uninstalled.png)

After the installation is complete, the user can choose to access the corresponding instance.

![instance installed page](./docs/images/dashboard-installed.png)

Click on PAUSE to release instance. After instance released, user can click on RECOVER to reload the AIGC instance, without losing user data.

![instance installed page](./docs/images/dashboard-recover.png)

Click on Logout, return to initial page.

![logout](./docs/images/logout.png)


## What's next

- [Install AIGC-Gateway](./docs/安装部署.md)

