
# Install RVM
command curl -sSL https://rvm.io/mpapis.asc | gpg --import
command curl -sSL https://rvm.io/pkuczynski.asc | gpg --import
curl -sSL https://get.rvm.io | bash -s stable

# Setup rvm envs
rvm requirements
rvm install 3.3.5
rvm install 3.3.4
rvm alias create default 3.3.5

rvm list

#
#